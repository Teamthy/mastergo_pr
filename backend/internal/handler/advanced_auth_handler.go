package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/models"
	"backend/internal/service"
)

type AdvancedAuthHandler struct {
	authService          *service.AuthService
	passwordResetService *service.PasswordResetService
	emailService         *service.EmailService
}

func NewAdvancedAuthHandler(
	authService *service.AuthService,
	passwordResetService *service.PasswordResetService,
	emailService *service.EmailService,
) *AdvancedAuthHandler {
	return &AdvancedAuthHandler{
		authService:          authService,
		passwordResetService: passwordResetService,
		emailService:         emailService,
	}
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
