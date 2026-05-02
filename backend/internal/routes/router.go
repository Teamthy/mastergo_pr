package routes

import (
	"backend/internal/handler"
	"backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	walletHandler *handler.WalletHandler,
	jwtSecret string,
) *chi.Mux {

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		AuthRoutes(r, authHandler, jwtSecret)

		// Wallet routes with proper JWT authentication
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtSecret))
			WalletRoutes(r, walletHandler)
		})
	})

	return r
}
