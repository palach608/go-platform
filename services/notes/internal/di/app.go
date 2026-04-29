package di

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/palach608/go-platform/pkg/auth"
	"github.com/palach608/go-platform/services/notes/internal/api"
	"github.com/palach608/go-platform/services/notes/internal/application"
	"github.com/palach608/go-platform/services/notes/internal/config"
	"github.com/palach608/go-platform/services/notes/internal/domain/service"
	"github.com/palach608/go-platform/services/notes/internal/infrastructure/db"
	"github.com/palach608/go-platform/services/notes/internal/infrastructure/repository"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(config.Load),
		fx.Provide(NewPool),

		fx.Provide(fx.Annotate(
			repository.NewPostgresNoteRepository,
			fx.As(new(repository.NoteRepository)),
		)),

		fx.Provide(fx.Annotate(
			application.NewNoteService,
			fx.As(new(service.NoteService)),
		)),

		fx.Provide(func(cfg *config.Config) *auth.Middleware {
			return auth.NewMiddleware(cfg.JWTSecret)
		}),

		fx.Provide(api.NewNoteHandler),
		fx.Provide(NewRouter),

		fx.Invoke(RunServer),
	)
}

func NewRouter(noteHandler *api.NoteHandler, middleware *auth.Middleware) chi.Router {
	return api.NewRouter(noteHandler, middleware.Auth)
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
