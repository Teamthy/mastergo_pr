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
