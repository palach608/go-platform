package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *Handler) chi.Router {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./api"))
	r.Handle("/api-docs/*", http.StripPrefix("/api-docs/", fs))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api-docs/swagger.yaml"),
	))

	r.Get("/api/v1/stats/top-users", handler.TopUsers)
	r.Get("/api/v1/stats/notes-per-day", handler.NotesPerDay)

	return r
}
