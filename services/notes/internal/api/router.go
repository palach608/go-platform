package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(noteHandler *NoteHandler, userHandler *UserHandler, middleware *Middleware) chi.Router {
	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("./web")))

	//Публичные маршруты без авторизации
	router.Get("/api/v1/notes", noteHandler.GetAll)            //Вывод всех заметок на главную
	router.Post("/api/v1/auth/register", userHandler.Register) //Возможность регистрации
	router.Post("/api/v1/auth/login", userHandler.Login)       //Возможность входа

	//Защищенные маршруты
	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Post("/api/v1/notes", noteHandler.Create)        //Создание заметки
		r.Get("/api/v1/notes/{id}", noteHandler.GetByID)   //Получение заметки по id
		r.Put("/api/v1/notes/{id}", noteHandler.Update)    //Обновление заметки
		r.Delete("/api/v1/notes/{id}", noteHandler.Delete) //Удаление заметки
	})

	return router
}
