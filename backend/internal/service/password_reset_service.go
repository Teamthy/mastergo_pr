package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

type PasswordResetService struct {
	db            *database.Database
	emailService  *EmailService
	tokenLifetime time.Duration
}

func NewPasswordResetService(db *database.Database, emailService *EmailService) *PasswordResetService {
	return &PasswordResetService{
		db:            db,
		emailService:  emailService,
		tokenLifetime: time.Hour, // 1 hour expiration
	}
}

// RequestPasswordReset creates a password reset token and sends email
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, email string, frontendURL string) error {
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not (security best practice)
		return nil
	}

	// Generate reset token
	token := s.generateToken()
	tokenHash := s.hashToken(token)

	resetToken := &models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.tokenLifetime),
		CreatedAt: time.Now(),
	}

	if err := s.db.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send reset email
	resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", frontendURL, token)
	return s.emailService.SendPasswordResetEmail(user.Email, resetLink)
}

// ValidateResetToken validates a reset token
func (s *PasswordResetService) ValidateResetToken(ctx context.Context, token string) (uuid.UUID, error) {
	tokenHash := s.hashToken(token)

	resetToken, err := s.db.GetPasswordResetTokenByHash(ctx, tokenHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid reset token: %w", err)
	}

	if resetToken.UsedAt != nil {
		return uuid.Nil, fmt.Errorf("reset token already used")
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return uuid.Nil, fmt.Errorf("reset token expired")
	}

	return resetToken.UserID, nil
}

// ResetPassword validates token and updates password
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	userID, err := s.ValidateResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Update password
	tokenHash := s.hashToken(token)
	now := time.Now()

	if err := s.db.UpdateUserPassword(ctx, userID, newPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.db.MarkPasswordResetTokenUsed(ctx, tokenHash, now); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

// Helper function to generate random token
func (s *PasswordResetService) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Helper function to hash token for storage
func (s *PasswordResetService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}
