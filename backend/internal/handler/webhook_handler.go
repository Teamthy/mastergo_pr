package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/models"
	"backend/internal/service"

	"github.com/google/uuid"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
	auditService   *service.AuditLogService
}

func NewWebhookHandler(webhookService *service.WebhookService, auditService *service.AuditLogService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
		auditService:   auditService,
	}
}

// CreateWebhook creates a new webhook
func (h *WebhookHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	var req models.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	webhook, err := h.webhookService.CreateWebhook(r.Context(), userID, req.URL, req.Events)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create webhook", err.Error())
		return
	}

	h.auditService.LogAction(r.Context(), &userID, "webhook.created", "webhook", webhook.ID.String(), nil, webhook, getIPAddress(r), r.UserAgent())

	respondJSON(w, http.StatusCreated, webhook)
}

// GetWebhooks retrieves all webhooks for the user
func (h *WebhookHandler) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	webhooks, err := h.webhookService.GetUserWebhooks(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve webhooks", err.Error())
		return
	}

	if webhooks == nil {
		webhooks = []models.Webhook{}
	}

	respondJSON(w, http.StatusOK, webhooks)
}

// UpdateWebhook updates a webhook
func (h *WebhookHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	webhookIDStr := r.PathValue("id")
	webhookID, err := uuid.Parse(webhookIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid webhook ID", err.Error())
		return
	}

	var req models.UpdateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	webhook, err := h.webhookService.UpdateWebhook(r.Context(), webhookID, req.URL, req.Events, req.Active)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update webhook", err.Error())
		return
	}

	h.auditService.LogAction(r.Context(), &userID, "webhook.updated", "webhook", webhook.ID.String(), nil, webhook, getIPAddress(r), r.UserAgent())

	respondJSON(w, http.StatusOK, webhook)
}

// DeleteWebhook deletes a webhook
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	webhookIDStr := r.PathValue("id")
	webhookID, err := uuid.Parse(webhookIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid webhook ID", err.Error())
		return
	}

	if err := h.webhookService.DeleteWebhook(r.Context(), webhookID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete webhook", err.Error())
		return
	}

	h.auditService.LogAction(r.Context(), &userID, "webhook.deleted", "webhook", webhookID.String(), nil, nil, getIPAddress(r), r.UserAgent())

	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Webhook deleted successfully",
	})
}

// Helper function to get IP address
func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
