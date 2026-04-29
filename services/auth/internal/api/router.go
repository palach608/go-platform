package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/palach608/go-platform/services/auth/internal/api/handlers"
)

func NewRouter(h *handlers.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/auth/register", h.Register)

	return r
}
