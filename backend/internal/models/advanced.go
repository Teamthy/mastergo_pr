package models

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"-"` // Never expose
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uuid.UUID              `json:"id"`
	UserID       *uuid.UUID             `json:"user_id,omitempty"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type,omitempty"`
	ResourceID   string                 `json:"resource_id,omitempty"`
	OldValues    map[string]interface{} `json:"old_values,omitempty"`
	NewValues    map[string]interface{} `json:"new_values,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	Status       string                 `json:"status"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID            uuid.UUID         `json:"id"`
	UserID        uuid.UUID         `json:"user_id"`
	URL           string            `json:"url"`
	Events        []string          `json:"events"`
	Headers       map[string]string `json:"headers,omitempty"`
	Active        bool              `json:"active"`
	RetryCount    int               `json:"retry_count"`
	LastTriggered *time.Time        `json:"last_triggered_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// WebhookEvent represents a webhook event
type WebhookEvent struct {
	ID           uuid.UUID              `json:"id"`
	WebhookID    uuid.UUID              `json:"webhook_id"`
	EventType    string                 `json:"event_type"`
	Payload      map[string]interface{} `json:"payload"`
	Status       string                 `json:"status"` // pending, sent, failed, retrying
	Attempts     int                    `json:"attempts"`
	MaxAttempts  int                    `json:"max_attempts"`
	ResponseCode *int                   `json:"response_code,omitempty"`
	ResponseBody string                 `json:"response_body,omitempty"`
	NextRetryAt  *time.Time             `json:"next_retry_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// APIAnalytics represents API usage analytics
type APIAnalytics struct {
	ID             uuid.UUID `json:"id"`
	APIKeyID       uuid.UUID `json:"api_key_id"`
	UserID         uuid.UUID `json:"user_id"`
	Endpoint       string    `json:"endpoint"`
	Method         string    `json:"method"`
	ResponseTimeMs int       `json:"response_time_ms"`
	StatusCode     int       `json:"status_code"`
	RequestSize    int       `json:"request_size,omitempty"`
	ResponseSize   int       `json:"response_size,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// RateLimitStats represents rate limit statistics
type RateLimitStats struct {
	ID           uuid.UUID `json:"id"`
	APIKeyID     uuid.UUID `json:"api_key_id"`
	Date         time.Time `json:"date"`
	RequestCount int       `json:"request_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
