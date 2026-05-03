package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/models"
	"backend/internal/service"
)

type AdvancedAuthHandler struct {
	authService          *service.AuthService
	twoFAService         *service.TwoFAService
	passwordResetService *service.PasswordResetService
	emailService         *service.EmailService
}

func NewAdvancedAuthHandler(
	authService *service.AuthService,
	twoFAService *service.TwoFAService,
	passwordResetService *service.PasswordResetService,
	emailService *service.EmailService,
) *AdvancedAuthHandler {
	return &AdvancedAuthHandler{
		authService:          authService,
		twoFAService:         twoFAService,
		passwordResetService: passwordResetService,
		emailService:         emailService,
	}
}

// Setup2FA initiates 2FA setup
func (h *AdvancedAuthHandler) Setup2FA(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	var req models.Setup2FARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// Verify password
	if _, err := h.authService.GetProfile(r.Context(), userID); err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "Invalid credentials")
		return
	}

	// TODO: Verify password against user's password hash

	// Generate 2FA setup
	setup, err := h.twoFAService.SetupTwoFA(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Setup failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, setup)
}

// Verify2FA verifies and enables 2FA
func (h *AdvancedAuthHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	var req models.Verify2FARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.Code == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "code is required")
		return
	}

	// Get the temporary 2FA setup
	setup, err := h.twoFAService.SetupTwoFA(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Verification failed", err.Error())
		return
	}

	// Enable 2FA
	if err := h.twoFAService.EnableTwoFA(r.Context(), userID, setup.Secret, req.Code, setup.BackupCodes); err != nil {
		respondError(w, http.StatusBadRequest, "Verification failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Two-factor authentication enabled successfully",
	})
}

// Disable2FA disables 2FA
func (h *AdvancedAuthHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	if err := h.twoFAService.DisableTwoFA(r.Context(), userID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to disable 2FA", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Two-factor authentication disabled successfully",
	})
}

// RequestPasswordReset sends password reset email
func (h *AdvancedAuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "email is required")
		return
	}

	// Get frontend URL from query params or header
	frontendURL := r.URL.Query().Get("frontend_url")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	if err := h.passwordResetService.RequestPasswordReset(r.Context(), req.Email, frontendURL); err != nil {
		respondError(w, http.StatusInternalServerError, "Request failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Password reset email sent successfully",
	})
}

// ResetPassword validates token and resets password
func (h *AdvancedAuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		respondError(w, http.StatusBadRequest, "Missing fields", "token and new_password are required")
		return
	}

	if err := h.passwordResetService.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		respondError(w, http.StatusBadRequest, "Password reset failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Password reset successfully",
	})
}
