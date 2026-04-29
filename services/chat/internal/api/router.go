package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/palach608/go-platform/pkg/auth"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(roomHandler *RoomHandler, wsHandler *WSHandler, middleware *auth.Middleware) chi.Router {
	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("./api"))
	router.Handle("/api-docs/*", http.StripPrefix("/api-docs/", fs))
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api-docs/swagger.yaml"),
	))

	router.Handle("/*", http.FileServer(http.Dir("./web")))

	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Post("/api/v1/rooms", roomHandler.Create)
		r.Get("/api/v1/rooms", roomHandler.GetAll)
		r.Get("/api/v1/rooms/{id}", roomHandler.GetByID)
		r.Delete("/api/v1/rooms/{id}", roomHandler.Delete)
		r.Get("/api/v1/rooms/{id}/messages", roomHandler.GetMessages)

		r.Get("/api/v1/ws/{roomID}", wsHandler.ServeWS)
	})

	return router
}
