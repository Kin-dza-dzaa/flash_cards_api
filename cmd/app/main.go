package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	_ "github.com/Kin-dza-dzaa/flash_cards_api/docs"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1/rest"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1/server"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/repository/googletrans"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/repository/postgresql"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/service"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/googletransclient"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"golang.org/x/exp/slog"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run(cfg))
}

func run(cfg config.Cfg) error {
	appCtx, cancelAppCtx := context.WithCancel(context.Background())
	defer cancelAppCtx()

	// Slog.
	l := logger.New(slog.Level(cfg.Logger.Level))

	// Jaeger.
	tp, err := otelTP(
		cfg.OpenTelemetry.ServiceName,
		cfg.OpenTelemetry.ServiceVersion,
		cfg.OpenTelemetry.Environment,
		cfg.OpenTelemetry.JaegerURL,
	)
	if err != nil {
		return fmt.Errorf("main - run - tracerProvider(): %w", err)
	}
	defer func(ctx context.Context) {
		shutdownCtx, cancel := context.WithTimeout(ctx, cfg.HTTP.ShutdownTimeout)
		defer cancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			l.Error(
				"couldn't shutdown tracer provider",
				slog.String("error", err.Error()),
			)
		}
	}(appCtx)
	otel.SetTracerProvider(tp)

	// Clients.
	client, err := googletransclient.New(cfg.GoogleAPI.URL)
	if err != nil {
		return fmt.Errorf("main - run - googletransclient.New: %w", err)
	}
	pool, err := postgres.New(appCtx, cfg.PG.URL, cfg.PG.MaxPoolSize)
	if err != nil {
		return fmt.Errorf("main - run - postgres.New: %w", err)
	}
	defer pool.Close()

	// Adapters/Repo layer.
	r := postgresql.NewWordPostgre(pool)
	g := googletrans.New(client, cfg.GoogleAPI.DefaultSrcLang, cfg.GoogleAPI.DefaultTrgtLang)

	// Usecase/business logic layer.
	s := service.NewWordService(r, g)

	// Port layer.
	h := rest.NewWordHandler(s, l)
	c := chi.NewRouter()
	h.Register(c, cfg)

	// Server start-up.
	srv := server.New(cfg, l, c)
	doneChan := srv.Start(appCtx)

	<-doneChan
	return nil
}

// Get Jaeger tracer provider.
func otelTP(serviceName, version, environment, url string) (*trace.TracerProvider, error) {
	// Create the Jaeger exporter.
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceVersion(version),
				semconv.ServiceName(serviceName),
				attribute.String("environment", environment),
			),
		),
	)
	return tp, nil
}
