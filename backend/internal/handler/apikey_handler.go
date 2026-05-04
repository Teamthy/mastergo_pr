package handler

import (
	"backend/internal/service"
	"encoding/json"
	"log"
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
		log.Println("ERROR: User context missing")
		writeError(w, http.StatusUnauthorized, "Unauthorized: User context missing")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusInternalServerError)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		log.Println("ERROR: API Key name is required")
		http.Error(w, "API Key name is required", http.StatusBadRequest)
		return
	}

	log.Printf("Creating API key for user %s with name %s", userID, req.Name)
	resp, err := h.service.CreateKey(r.Context(), userID, req.Name)
	if err != nil {
		log.Printf("ERROR: Failed to create API key: %v", err)
		http.Error(w, "Failed to create API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("API key created successfully for user %s", userID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ApiKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	log.Printf("Listing API keys for user %s", userID)
	keys, err := h.service.ListKeys(r.Context(), userID)
	if err != nil {
		log.Printf("ERROR: Failed to list API keys: %v", err)
		http.Error(w, "Failed to fetch keys", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d API keys for user %s", len(keys), userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

func (h *ApiKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	keyIDStr := chi.URLParam(r, "id")
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		log.Printf("ERROR: Invalid API Key ID: %v", err)
		http.Error(w, "Invalid API Key ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Context().Value("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	log.Printf("Deleting API key %s for user %s", keyID, userID)
	if err := h.service.RevokeKey(r.Context(), keyID, userID); err != nil {
		log.Printf("ERROR: Failed to revoke key: %v", err)
		http.Error(w, "Failed to revoke key", http.StatusInternalServerError)
		return
	}

	log.Printf("API key %s deleted successfully", keyID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *ApiKeyHandler) Regenerate(w http.ResponseWriter, r *http.Request) {
	keyIDStr := chi.URLParam(r, "id")
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		log.Printf("ERROR: Invalid API Key ID: %v", err)
		http.Error(w, "Invalid API Key ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Context().Value("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	log.Printf("Regenerating API key %s for user %s", keyID, userID)
	resp, err := h.service.RegenerateKey(r.Context(), keyID, userID)
	if err != nil {
		log.Printf("ERROR: Failed to regenerate key: %v", err)
		http.Error(w, "Failed to regenerate key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("API key %s regenerated successfully", keyID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
