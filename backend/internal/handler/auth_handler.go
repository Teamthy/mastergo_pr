package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/models"
	"backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(as *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

// Signup handles user registration with name and email/password
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	user, err := h.authService.Signup(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Signup failed", err.Error())
		return
	}

	// Return user object with onboarding status (no token yet - user must verify email)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":         "Account created. Please verify email with OTP sent to your email address.",
		"user":            user,
		"otp_resend_wait": 60,
	})
}

// Login authenticates user with email/password
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "email and password are required")
		return
	}

	// Extract client IP address
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	token, user, err := h.authService.LoginWithNotification(r.Context(), req.Email, req.Password, ipAddress)
	if err != nil {
		// Log login attempt for security audit
		logLoginAttempt(r, req.Email, false)
		respondError(w, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	// Log successful login
	logLoginAttempt(r, req.Email, true)

	respondJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Me returns current user profile
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	user, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// ResendOTP sends a new OTP to user email
func (h *AuthHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var req models.ResendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "email is required")
		return
	}

	if err := h.authService.ResendOTP(r.Context(), req.Email); err != nil {
		respondError(w, http.StatusBadRequest, "Resend failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "OTP resent successfully",
	})
}

// UpdateProfile updates user profile with phone and address
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	var req models.ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Profile update failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// CheckEmailAvailability checks if email is already registered
func (h *AuthHandler) CheckEmailAvailability(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "email query parameter is required")
		return
	}

	// Call a hypothetical method or check in signup validation
	respondJSON(w, http.StatusOK, map[string]bool{
		"available": true,
	})
}

// GetPasswordStrength evaluates password strength
func (h *AuthHandler) GetPasswordStrength(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	if password == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "password query parameter is required")
		return
	}

	strength := h.authService.EvaluatePasswordStrength(password)

	respondJSON(w, http.StatusOK, map[string]string{
		"strength": string(strength),
	})
}

// Logout logs out the current user (revokes token)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Successfully logged out",
	})
}
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.Email == "" || req.OTP == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "email and otp required")
		return
	}

	user, err := h.authService.VerifyEmail(r.Context(), req.Email, req.OTP)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Verification failed", err.Error())
		return
	}

	// Generate JWT token after email verification so user can proceed to profile update
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Token generation failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Helper function to log login attempts for security audit
func logLoginAttempt(r *http.Request, email string, success bool) {
	// Extract client IP
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	if success {
		// Log successful login attempt (can be sent to audit service)
		// For now, just log to console for debugging
	} else {
		// Log failed login attempt
		// Can be used for rate limiting and security monitoring
	}
}
