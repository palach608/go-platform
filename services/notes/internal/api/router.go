package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(noteHandler *NoteHandler, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("./api"))
	router.Handle("/api-docs/*", http.StripPrefix("/api-docs/", fs))
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api-docs/swagger.yaml"),
	))

	router.Get("/api/v1/notes", noteHandler.GetAll)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/api/v1/notes", noteHandler.Create)
		r.Get("/api/v1/notes/{id}", noteHandler.GetByID)
		r.Put("/api/v1/notes/{id}", noteHandler.Update)
		r.Delete("/api/v1/notes/{id}", noteHandler.Delete)
	})

	return router
}
