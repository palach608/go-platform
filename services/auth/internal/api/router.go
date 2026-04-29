package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/palach608/go-platform/services/auth/internal/api/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handlers.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Put("/user", h.Update)
		r.Delete("/user/{id}", h.Delete)
	})

	fs := http.FileServer(http.Dir("./api"))
	r.Handle("/api/*", http.StripPrefix("/api/", fs))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api/swagger.yaml"),
	))

	return r
}
