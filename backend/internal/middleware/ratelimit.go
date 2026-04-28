package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimitConfig defines rate limiting rules
type RateLimitConfig struct {
	Requests int           // number of requests
	Window   time.Duration // time window
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rdb *redis.Client, limits map[string]RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get endpoint path
			path := r.URL.Path
			method := r.Method
			endpoint := fmt.Sprintf("%s %s", method, path)

			// Get client IP
			ip := getClientIP(r)
			key := fmt.Sprintf("ratelimit:%s:%s", endpoint, ip)

			// Check if endpoint has rate limit
			limit, exists := limits[endpoint]
			if !exists {
				next.ServeHTTP(w, r)
				return
			}

			// Check current count in Redis
			ctx := r.Context()
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				// Log error but allow request to proceed
				next.ServeHTTP(w, r)
				return
			}

			// Set expiration on first increment
			if count == 1 {
				rdb.Expire(ctx, key, limit.Window)
			}

			// Check if limit exceeded
			if count > int64(limit.Requests) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(limit.Window.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{"error":"rate_limit_exceeded","message":"Too many requests. Please try again later."}`)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header (proxy)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// ApplyRateLimits applies rate limiting to auth routes
func ApplyRateLimits(rdb *redis.Client) func(http.Handler) http.Handler {
	limits := map[string]RateLimitConfig{
		"POST /auth/signup":       {Requests: 5, Window: 1 * time.Minute},  // 5 per minute
		"POST /auth/login":        {Requests: 10, Window: 1 * time.Minute}, // 10 per minute
		"POST /auth/verify-email": {Requests: 5, Window: 1 * time.Minute},  // 5 per minute
		"POST /auth/resend-otp":   {Requests: 3, Window: 1 * time.Minute},  // 3 per minute
	}

	return RateLimitMiddleware(rdb, limits)
}
