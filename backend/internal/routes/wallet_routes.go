package routes

import (
	"backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func WalletRoutes(r chi.Router, h *handler.WalletHandler) {
	r.Route("/wallet", func(r chi.Router) {
		r.Post("/create", h.Create)
	})
}
