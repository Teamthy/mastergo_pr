package models

import (
	"time"

	"github.com/google/uuid"
)

// Wallet represents an Ethereum wallet for a user
type Wallet struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	UserID              uuid.UUID `json:"user_id" db:"user_id"`
	Address             string    `json:"address" db:"address"`
	EncryptedPrivateKey string    `json:"-" db:"encrypted_private_key"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// Balance represents the ledger balance for a user
type Balance struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	BalanceWei string    `json:"balance_wei" db:"balance"` // Stored as decimal/numeric
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction represents a transaction in the internal ledger
type Transaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TxHash    string    `json:"tx_hash" db:"tx_hash"`
	Type      string    `json:"type" db:"type"` // "deposit" or "withdrawal"
	AmountWei string    `json:"amount_wei" db:"amount_wei"`
	To        string    `json:"to" db:"to"`
	Status    string    `json:"status" db:"status"` // "pending", "confirmed", "failed"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// WithdrawRequest is the request payload for withdrawals
type WithdrawRequest struct {
	AmountWei string `json:"amount_wei" validate:"required"`
	To        string `json:"to" validate:"required"`
}
