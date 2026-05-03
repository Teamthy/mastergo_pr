package handler

import (
	"net/http"

	"github.com/google/uuid"
)

func getUserID(r *http.Request) (uuid.UUID, error) {
	val := r.Context().Value("user_id")
	idStr, ok := val.(string)
	if !ok || idStr == "" {
		return uuid.Nil, http.ErrNoCookie
	}
	return uuid.Parse(idStr)
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
