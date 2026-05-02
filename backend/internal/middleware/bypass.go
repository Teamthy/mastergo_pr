package middleware

import (
	"context"
	"net/http"
)

// BypassAuthMiddleware sets a mock user ID for development/testing without requiring authentication
func BypassAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use a fixed UUID for testing - you can change this to any valid UUID
		mockUserID := "550e8400-e29b-41d4-a716-446655440000"
		ctx := context.WithValue(r.Context(), "user_id", mockUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
