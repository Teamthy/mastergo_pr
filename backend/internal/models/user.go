package models

import (
	"time"

	"github.com/google/uuid"
)

type OnboardingStep string

const (
	StepStart            OnboardingStep = "START"
	StepEmailEntered     OnboardingStep = "EMAIL_ENTERED"
	StepEmailVerified    OnboardingStep = "EMAIL_VERIFIED"
	StepProfileCompleted OnboardingStep = "PROFILE_COMPLETED"
	StepCompleted        OnboardingStep = "COMPLETED"
)

type User struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Email            string         `json:"email" db:"email"`
	PasswordHash     string         `json:"-" db:"password_hash"`
	FirstName        string         `json:"first_name" db:"first_name"`
	LastName         string         `json:"last_name" db:"last_name"`
	Phone            *string        `json:"phone" db:"phone"`
	Address          *string        `json:"address" db:"address"`
	OnboardingStatus OnboardingStep `json:"onboarding_status" db:"onboarding_status"`
	EmailVerified    bool           `json:"email_verified" db:"email_verified"`
	LastLoginAt      *time.Time     `json:"last_login_at" db:"last_login_at"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

func NewUser(firstName, lastName, email, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:               uuid.New(),
		FirstName:        firstName,
		LastName:         lastName,
		Email:            email,
		PasswordHash:     passwordHash,
		OnboardingStatus: StepEmailEntered,
		EmailVerified:    false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
