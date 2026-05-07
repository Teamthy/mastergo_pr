package service

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	repo         *repository.UserRepository
	rdb          *redis.Client
	jwtSecret    string
	emailService *EmailService
}

type PasswordStrength string

const (
	PasswordWeak   PasswordStrength = "weak"
	PasswordMedium PasswordStrength = "medium"
	PasswordStrong PasswordStrength = "strong"
)

func NewAuthService(repo *repository.UserRepository, rdb *redis.Client, jwtSecret string, emailService *EmailService) *AuthService {
	return &AuthService{
		repo:         repo,
		rdb:          rdb,
		jwtSecret:    jwtSecret,
		emailService: emailService,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func otpKey(email string) string {
	return "otp:" + normalizeEmail(email)
}

func generateOTP() (string, error) {
	n, err := crand.Int(crand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// ValidatePassword checks password meets industry standards
func (s *AuthService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return fmt.Errorf("password must contain at least one number")
	}

	hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSymbol {
		return fmt.Errorf("password must contain at least one special character")
	}

	if regexp.MustCompile(`\s`).MatchString(password) {
		return fmt.Errorf("password must not contain spaces")
	}

	return nil
}

// EvaluatePasswordStrength returns password strength level
func (s *AuthService) EvaluatePasswordStrength(password string) PasswordStrength {
	score := 0

	if len(password) >= 12 {
		score++
	}
	if len(password) >= 16 {
		score++
	}

	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		score++
	}
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		score++
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		score++
	}
	if regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		score++
	}

	switch {
	case score <= 2:
		return PasswordWeak
	case score <= 4:
		return PasswordMedium
	default:
		return PasswordStrong
	}
}

func (s *AuthService) Signup(ctx context.Context, req *models.SignUpRequest) (*models.User, error) {
	// Normalize email
	req.Email = normalizeEmail(req.Email)

	// Validate names
	if strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.LastName) == "" {
		return nil, fmt.Errorf("first name and last name are required")
	}

	if len(req.FirstName) < 2 || len(req.FirstName) > 50 {
		return nil, fmt.Errorf("first name must be 2-50 characters")
	}

	if len(req.LastName) < 2 || len(req.LastName) > 50 {
		return nil, fmt.Errorf("last name must be 2-50 characters")
	}

	// Validate name contains only letters
	if !regexp.MustCompile(`^[a-zA-Z\s'-]+$`).MatchString(req.FirstName) {
		return nil, fmt.Errorf("first name must contain only letters")
	}

	if !regexp.MustCompile(`^[a-zA-Z\s'-]+$`).MatchString(req.LastName) {
		return nil, fmt.Errorf("last name must contain only letters")
	}

	// Validate email
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Validate password matches
	if req.Password != req.ConfirmPassword {
		return nil, fmt.Errorf("passwords do not match")
	}

	// Validate password strength
	if err := s.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check email uniqueness
	existing, _ := s.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:               uuid.New(),
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PasswordHash:     string(hashedPassword),
		OnboardingStatus: models.StepEmailEntered,
		EmailVerified:    false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	user, err = s.repo.CreateUser(
		ctx,
		req.FirstName,
		req.LastName,
		req.Email,
		string(hashedPassword),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, err
	}

	key := otpKey(req.Email)
	if err := s.rdb.Set(ctx, key, otp, 5*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("failed to store otp: %w", err)
	}

	// Store attempt counter
	attemptsKey := "otp_attempts:" + req.Email
	s.rdb.Set(ctx, attemptsKey, 0, 5*time.Minute)

	// Send OTP via email asynchronously (non-blocking)
	if s.emailService != nil {
		s.emailService.SendOTPEmailAsync(req.Email, otp)
	}

	log.Printf("OTP DEBUG KEY=%s OTP=%s", key, otp)
	fmt.Printf("\n--- [MOCK EMAIL FALLBACK] ---\nTo: %s\nOTP: %s\n--------------------\n\n", req.Email, otp)

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	email = normalizeEmail(email)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if !user.EmailVerified {
		return "", nil, fmt.Errorf("please verify your email first")
	}

	token, err := utils.GenerateToken(user.ID.String(), s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.Update(ctx, user); err != nil {
		log.Printf("warning: failed to update last login: %v", err)
	}

	return token, user, nil
}

// LoginWithNotification performs login and sends notification email
func (s *AuthService) LoginWithNotification(ctx context.Context, email, password, ipAddress string) (string, *models.User, error) {
	token, user, err := s.Login(ctx, email, password)
	if err != nil {
		return "", nil, err
	}

	// Send login notification email
	if s.emailService != nil && user != nil {
		userName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		if err := s.emailService.SendLoginNotificationEmail(email, userName, ipAddress); err != nil {
			log.Printf("Failed to send login notification email: %v", err)
			// Don't fail login if email sending fails
		}
	}

	return token, user, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) (*models.User, error) {
	email = normalizeEmail(email)
	code = strings.TrimSpace(code)

	if email == "" || code == "" {
		return nil, fmt.Errorf("email and otp are required")
	}

	key := otpKey(email)

	stored, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("otp expired or not found")
	}

	stored = strings.TrimSpace(stored)

	log.Printf("OTP VERIFY DEBUG key=%s stored=%q input=%q", key, stored, code)

	if stored != code {
		// Increment attempt counter
		attemptsKey := "otp_attempts:" + email
		attempts, _ := s.rdb.Incr(ctx, attemptsKey).Result()

		if attempts >= 5 {
			s.rdb.Del(ctx, key)
			return nil, fmt.Errorf("too many OTP attempts, please request a new one")
		}

		return nil, fmt.Errorf("invalid OTP")
	}

	if err := s.repo.MarkUserVerified(ctx, email); err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	s.rdb.Del(ctx, key)
	s.rdb.Del(ctx, "otp_attempts:"+email)

	// Return the verified user object
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Send welcome email after successful verification
	if s.emailService != nil {
		userName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		if err := s.emailService.SendWelcomeEmail(email, userName); err != nil {
			log.Printf("Failed to send welcome email: %v", err)
			// Don't fail verification if welcome email fails
		}
	}

	return user, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *AuthService) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	req *models.ProfileUpdateRequest,
) (*models.User, error) {
	if req == nil {
		return nil, fmt.Errorf("profile request is required")
	}

	phone := strings.TrimSpace(req.Phone)
	address := strings.TrimSpace(req.Address)

	if phone == "" || address == "" {
		return nil, fmt.Errorf("phone and address are required")
	}

	if len(phone) < 10 || len(phone) > 20 {
		return nil, fmt.Errorf("phone must be 10-20 characters")
	}

	if len(address) < 10 {
		return nil, fmt.Errorf("address must be at least 10 characters")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.Phone = &phone
	user.Address = &address
	user.OnboardingStatus = models.StepCompleted // Mark onboarding as completed
	user.UpdatedAt = time.Now()

	if err := s.repo.UpdateProfile(ctx, userID, user.FirstName, user.LastName, phone, address, models.StepCompleted); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return user, nil
}

func (s *AuthService) ResendOTP(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	if email == "" {
		return fmt.Errorf("email is required")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("invalid email")
	}

	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	cooldownKey := "otp:cooldown:" + email
	ok, err := s.rdb.Exists(ctx, cooldownKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check otp cooldown: %w", err)
	}
	if ok == 1 {
		return fmt.Errorf("please wait before requesting another otp")
	}

	otp, err := generateOTP()
	if err != nil {
		return err
	}

	key := otpKey(email)
	if err := s.rdb.Set(ctx, key, otp, 10*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to store otp: %w", err)
	}

	if err := s.rdb.Set(ctx, cooldownKey, "1", 60*time.Second).Err(); err != nil {
		return fmt.Errorf("failed to set otp cooldown: %w", err)
	}

	// Send OTP via email asynchronously (non-blocking)
	if s.emailService != nil {
		s.emailService.SendOTPEmailAsync(email, otp)
	}

	log.Printf("OTP RESEND DEBUG KEY=%s", key)
	fmt.Printf("\n--- [MOCK EMAIL FALLBACK] ---\nTo: %s\nOTP: %s\n--------------------\n\n", email, otp)

	return nil
}

// GenerateToken creates JWT token for authenticated user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	return utils.GenerateToken(user.ID.String(), s.jwtSecret)
}
