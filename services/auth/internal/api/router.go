package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/palach608/go-platform/services/auth/internal/api/handlers"
)

func NewRouter(h *handlers.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Put("/user", h.Update)
		r.Delete("/user/{id}", h.Delete)
	})

	return r
}
