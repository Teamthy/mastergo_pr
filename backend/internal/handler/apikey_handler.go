package handler

import (
	"backend/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ApiKeyHandler struct {
	service *service.ApiKeyService
}

func NewApiKeyHandler(service *service.ApiKeyService) *ApiKeyHandler {
	return &ApiKeyHandler{service: service}
}

func (h *ApiKeyHandler) Create(w http.ResponseWriter, r *http.Request) {

	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized: User context missing")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusInternalServerError)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "API Key name is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateKey(r.Context(), userID, req.Name)
	if err != nil {
		http.Error(w, "Failed to create API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ApiKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	keys, err := h.service.ListKeys(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

func (h *ApiKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	keyIDStr := chi.URLParam(r, "id")
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		http.Error(w, "Invalid API Key ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Context().Value("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	if err := h.service.RevokeKey(r.Context(), keyID, userID); err != nil {
		http.Error(w, "Failed to revoke key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
