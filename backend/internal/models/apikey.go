package models

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Name         string     `json:"name" db:"name"`
	PublicKey    string     `json:"public_key" db:"public_key"`
	HashedSecret string     `json:"-" db:"hashed_secret"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

type CreateApiKeyRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type ApiKeyCreateResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	SecretKey string    `json:"secret_key"`
	CreatedAt time.Time `json:"created_at"`
}
