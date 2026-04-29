package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *UserHandler, roomHandler *RoomHandler, wsHandler *WSHandler, middleware *Middleware) chi.Router {
	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("./web")))

	router.Post("/api/v1/auth/register", userHandler.Register)
	router.Post("/api/v1/auth/login", userHandler.Login)

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
