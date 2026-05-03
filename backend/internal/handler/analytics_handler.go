package handler

import (
	"net/http"
	"strconv"

	"backend/internal/service"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetAPIAnalytics retrieves API analytics for an API key
func (h *AnalyticsHandler) GetAPIAnalytics(w http.ResponseWriter, r *http.Request) {
	_, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	apiKeyID := r.PathValue("id")
	if apiKeyID == "" {
		respondError(w, http.StatusBadRequest, "Missing API key ID", "")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	// Parse API key ID to UUID
	apiKeyUUID, err := parseUUID(apiKeyID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid API key ID", err.Error())
		return
	}

	analytics, err := h.analyticsService.GetAPIAnalytics(r.Context(), apiKeyUUID, days)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve analytics", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, analytics)
}

// GetUserAnalytics retrieves analytics for all APIs of a user
func (h *AnalyticsHandler) GetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	analytics, err := h.analyticsService.GetUserAnalytics(r.Context(), userID, days)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve analytics", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, analytics)
}

// ExportAnalytics exports analytics as CSV
func (h *AnalyticsHandler) ExportAnalytics(w http.ResponseWriter, r *http.Request) {
	_, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}

	apiKeyID := r.PathValue("id")
	if apiKeyID == "" {
		respondError(w, http.StatusBadRequest, "Missing API key ID", "")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	apiKeyUUID, err := parseUUID(apiKeyID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid API key ID", err.Error())
		return
	}

	csv, err := h.analyticsService.ExportAnalytics(r.Context(), apiKeyUUID, days)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to export analytics", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=analytics.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}
