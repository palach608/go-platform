package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/palach608/go-platform/pkg/auth"
)

func NewRouter(noteHandler *NoteHandler, jwtSecret string) *chi.Mux {
	router := chi.NewRouter()

	// Публичные маршруты (без авторизации)
	router.Get("/api/v1/notes", noteHandler.GetAll)

	// Защищённые маршруты
	router.Group(func(r chi.Router) {
		r.Use(auth.NewMiddleware(jwtSecret).Auth)

		r.Post("/api/v1/notes", noteHandler.Create)
		r.Get("/api/v1/notes/{id}", noteHandler.GetByID)
		r.Put("/api/v1/notes/{id}", noteHandler.Update)
		r.Delete("/api/v1/notes/{id}", noteHandler.Delete)
	})

	return router
}
