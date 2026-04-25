package service

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"log"
	"math/big"
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
	repo      *repository.UserRepository
	rdb       *redis.Client
	jwtSecret string
}

func NewAuthService(repo *repository.UserRepository, rdb *redis.Client, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		rdb:       rdb,
		jwtSecret: jwtSecret,
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

func (s *AuthService) Signup(ctx context.Context, email, password string) (*models.User, error) {
	email = normalizeEmail(email)

	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:            email,
		PasswordHash:     string(hashedPassword),
		OnboardingStatus: models.StepStart,
		IsVerified:       false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, err
	}

	key := otpKey(email)
	if err := s.rdb.Set(ctx, key, otp, 10*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("failed to store otp: %w", err)
	}

	log.Printf("OTP DEBUG KEY=%s", key)
	fmt.Printf("\n--- [MOCK EMAIL] ---\nTo: %s\nOTP: %s\n--------------------\n\n", email, otp)

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

	if !user.IsVerified {
		return "", nil, fmt.Errorf("please verify your email first")
	}

	token, err := utils.GenerateToken(user.ID.String(), s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) error {
	email = normalizeEmail(email)
	code = strings.TrimSpace(code)

	key := otpKey(email)

	stored, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("otp expired or not found")
	}

	stored = strings.TrimSpace(stored)

	log.Printf("OTP VERIFY DEBUG key=%s stored=%q input=%q", key, stored, code)

	if stored != code {
		return fmt.Errorf("invalid OTP")
	}

	if err := s.repo.MarkUserVerified(ctx, email); err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	if err := s.rdb.Del(ctx, key).Err(); err != nil {
		log.Printf("warning: failed to delete otp key %s: %v", key, err)
	}

	return nil
}
func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *AuthService) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	fullName string,
) (*models.User, error) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return nil, fmt.Errorf("full_name is required")
	}

	if err := s.repo.UpdateProfile(ctx, userID, fullName, models.StepProfileCompleted); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, userID)
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

	if user.IsVerified {
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

	log.Printf("OTP RESEND DEBUG KEY=%s", key)
	fmt.Printf("\n--- [MOCK EMAIL] ---\nTo: %s\nOTP: %s\n--------------------\n\n", email, otp)

	return nil
}
