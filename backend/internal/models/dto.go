package models

// SignUp Request
type SignUpRequest struct {
	FirstName       string `json:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName        string `json:"last_name" validate:"required,min=2,max=50,alpha"`
	Email           string `json:"email" validate:"required,email,lowercase"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// Login Request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Verify Email Request
type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6,numeric"`
}

// Resend OTP Request
type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Profile Update Request
type ProfileUpdateRequest struct {
	Phone   string `json:"phone" validate:"required,min=10,max=20"`
	Address string `json:"address" validate:"required,min=10"`
}

// Auth Response
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// OTP Response
type OTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

// Generic Success Response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error Response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Password Reset Request
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Password Reset Confirm Request
type PasswordResetConfirmRequest struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// Account Recovery Request
type AccountRecoveryRequest struct {
	Email            string `json:"email" validate:"required,email"`
	RecoveryEmail    string `json:"recovery_email" validate:"required,email"`
	SecurityQuestion string `json:"security_question"`
	SecurityAnswer   string `json:"security_answer"`
}

// Webhook Request
type CreateWebhookRequest struct {
	URL    string   `json:"url" validate:"required,url"`
	Events []string `json:"events" validate:"required,min=1"`
}

// Webhook Update Request
type UpdateWebhookRequest struct {
	URL    string   `json:"url" validate:"required,url"`
	Events []string `json:"events" validate:"required,min=1"`
	Active bool     `json:"active"`
}

// Transaction Filter Request
type TransactionFilterRequest struct {
	Type      string `json:"type"`   // "deposit" or "withdrawal"
	Status    string `json:"status"` // "pending", "confirmed", "failed"
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// Analytics Response
type AnalyticsResponse struct {
	TotalRequests   int64      `json:"total_requests"`
	RequestsToday   int        `json:"requests_today"`
	AvgResponseTime float64    `json:"avg_response_time"`
	TopEndpoints    []Endpoint `json:"top_endpoints"`
	ErrorRate       float64    `json:"error_rate"`
}

type Endpoint struct {
	Path       string  `json:"path"`
	Calls      int     `json:"calls"`
	AvgTime    float64 `json:"avg_time"`
	ErrorCount int     `json:"error_count"`
}
