package di

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/palach608/go-platform/services/auth/internal/api"
	"github.com/palach608/go-platform/services/auth/internal/api/handlers"
	"github.com/palach608/go-platform/services/auth/internal/application"
	"github.com/palach608/go-platform/services/auth/internal/config"
	"github.com/palach608/go-platform/services/auth/internal/infrastructure/repository"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.Load,
			NewDatabase,
			repository.NewPostgresRepo,
			application.NewAuthService,
			handlers.NewAuthHandler,
			api.NewRouter,
		),
		fx.Invoke(RunServer),
	)
}

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
}

func RunServer(lifecycle fx.Lifecycle, router *chi.Mux, cfg *config.Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Printf("Auth Service launch → http://localhost:%s\n", cfg.HTTPPort)
			go http.ListenAndServe(":"+cfg.HTTPPort, router)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Auth Service...")
			return nil
		},
	})
}
