package di

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	kafkapkg "github.com/palach608/go-platform/pkg/broker"
	"github.com/palach608/go-platform/services/analytics/internal/api"
	chclient "github.com/palach608/go-platform/services/analytics/internal/clickhouse"
	"github.com/palach608/go-platform/services/analytics/internal/config"
	"github.com/palach608/go-platform/services/analytics/internal/consumer"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(config.Load),
		fx.Provide(NewClickhouseClient),
		fx.Provide(consumer.NewEventHandler),
		fx.Provide(api.NewHandler),
		fx.Provide(api.NewRouter),
		fx.Invoke(RunServer),
	)
}

func NewClickhouseClient(cfg *config.Config) (*chclient.Client, error) {
	return chclient.NewClient(cfg.ClickhouseDSN)
}

func RunServer(router chi.Router, cfg *config.Config, handler *consumer.EventHandler) {
	go func() {
		brokers := strings.Split(cfg.KafkaBrokers, ",")
		c := kafkapkg.NewConsumer(brokers, kafkapkg.TopicNoteCreated, "analytics-group")
		c.Consume(context.Background(), handler.Handle)
	}()

	fmt.Printf("Analytics Service launch → http://localhost:%s\n", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		panic(err)
	}
}
