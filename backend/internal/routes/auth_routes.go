package routes

import (
	"backend/internal/handler"
	"backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, h *handler.AuthHandler, jwtSecret string) {
	r.Route("/auth", func(r chi.Router) {
		// Public routes
		r.Post("/signup", h.Signup)
		r.Post("/login", h.Login)
		r.Post("/verify-email", h.VerifyEmail)
		r.Post("/resend-otp", h.ResendOTP)
		r.Get("/email-available", h.CheckEmailAvailability)
		r.Get("/password-strength", h.GetPasswordStrength)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtSecret))
			r.Get("/me", h.Me)
			r.Patch("/profile", h.UpdateProfile)
			r.Post("/logout", h.Logout)
		})
	})
}
