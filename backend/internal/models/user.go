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
	FullName         string         `json:"full_name" db:"full_name"`
	OnboardingStatus OnboardingStep `json:"onboarding_status" db:"onboarding_status"`
	IsVerified       bool           `json:"is_verified" db:"is_verified"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
}

func NewUser(email, passwordHash string) *User {
	return &User{
		ID:               uuid.New(),
		Email:            email,
		PasswordHash:     passwordHash,
		OnboardingStatus: StepStart,
		IsVerified:       false,
		CreatedAt:        time.Now(),
	}
}
