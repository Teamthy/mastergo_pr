package handler

import (
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/service"
)

type AuditLogHandler struct {
	auditService *service.AuditLogService
}

func NewAuditLogHandler(auditService *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{auditService: auditService}
}

// GetAuditLogs retrieves audit logs for the user
func (h *AuditLogHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	logs, total, err := h.auditService.GetAuditLogs(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve audit logs", err.Error())
		return
	}

	if logs == nil {
		logs = []models.AuditLog{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"logs":  logs,
		"total": total,
	})
}

// GetAuditLogsSummary retrieves summary of audit logs
func (h *AuditLogHandler) GetAuditLogsSummary(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	summary, err := h.auditService.GetActionSummary(r.Context(), userID, days)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve summary", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, summary)
}
