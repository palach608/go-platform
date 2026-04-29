package di

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/palach608/go-platform/pkg/auth" // ← добавь этот импорт
	"github.com/palach608/go-platform/services/auth/internal/config"
	"github.com/palach608/go-platform/services/notes/internal/api"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(config.Load),
		fx.Provide(NewPool),
		fx.Provide(repository.NewPostgresNoteRepository), // UserRepository пока убираем
		fx.Provide(application.NewNoteService),

		// Новый middleware из pkg
		fx.Provide(NewJWTMiddleware),

		fx.Provide(api.NewNoteHandler),
		fx.Provide(api.NewRouter), // ← будет принимать jwtSecret
		fx.Invoke(RunServer),
	)
}

// Новый middleware из pkg/auth
func NewJWTMiddleware(cfg *config.Config) *auth.Middleware {
	return auth.NewMiddleware(cfg.JWTSecret)
}

// Обновлённый NewRouter — принимает только noteHandler и jwtSecret
func NewRouter(noteHandler *api.NoteHandler, middleware *auth.Middleware) chi.Router {
	return api.NewRouter(noteHandler, middleware) // пока так, потом можно упростить
}

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	return db.NewPool(cfg.DBDSN)
}

func RunServer(router chi.Router, cfg *config.Config) {
	if err := runMigrations(cfg); err != nil {
		panic(err)
	}
	fmt.Printf("Notes Service launch → http://localhost:%s\n", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		panic(err)
	}
}

func runMigrations(cfg *config.Config) error {
	m, err := migrate.New("file://migrations", cfg.DBDSN)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
