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

	goredis "github.com/redis/go-redis/v9"

	"github.com/palach608/go-platform/pkg/auth"
	"github.com/palach608/go-platform/services/chat/internal/api"
	"github.com/palach608/go-platform/services/chat/internal/application"
	"github.com/palach608/go-platform/services/chat/internal/config"
	"github.com/palach608/go-platform/services/chat/internal/domain/service"
	"github.com/palach608/go-platform/services/chat/internal/hub"
	"github.com/palach608/go-platform/services/chat/internal/infrastructure/db"
	redisClient "github.com/palach608/go-platform/services/chat/internal/infrastructure/redis"
	"github.com/palach608/go-platform/services/chat/internal/infrastructure/repository"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(config.Load),
		fx.Provide(NewPool),
		fx.Provide(NewRedisClient),

		fx.Provide(fx.Annotate(
			repository.NewPostgresRoomRepository,
			fx.As(new(repository.RoomRepository)),
		)),
		fx.Provide(fx.Annotate(
			repository.NewPostgresMessageRepository,
			fx.As(new(repository.MessageRepository)),
		)),

		fx.Provide(fx.Annotate(
			application.NewRoomService,
			fx.As(new(service.RoomService)),
		)),
		fx.Provide(fx.Annotate(
			application.NewMessageService,
			fx.As(new(service.MessageService)),
		)),

		fx.Provide(hub.NewHub),

		fx.Provide(func(cfg *config.Config) *auth.Middleware {
			return auth.NewMiddleware(cfg.JWTSecret)
		}),

		fx.Provide(api.NewRoomHandler),
		fx.Provide(api.NewWSHandler),
		fx.Provide(api.NewRouter),

		fx.Invoke(RunServer),
	)
}

func RunServer(router chi.Router, cfg *config.Config, h *hub.Hub) {
	if err := runMigrations(cfg); err != nil {
		panic(err)
	}
	go h.Run()
	fmt.Printf("Chat Service launch → http://localhost:%s\n", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		panic(err)
	}
}

func runMigrations(cfg *config.Config) error {
	m, err := migrate.New(
		"file://migrations",
		cfg.DBDSN,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	return db.NewPool(cfg.DBDSN)
}

func NewRedisClient(cfg *config.Config) (*goredis.Client, error) {
	return redisClient.NewClient(cfg.RedisDSN)
}
