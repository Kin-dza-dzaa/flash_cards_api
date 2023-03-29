// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	v1 "github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1"
	googletrans "github.com/Kin-dza-dzaa/flash_cards_api/internal/repository/google_translate"
	wordpostgres "github.com/Kin-dza-dzaa/flash_cards_api/internal/repository/word_postgres"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/usecase"
	googletransclient "github.com/Kin-dza-dzaa/flash_cards_api/pkg/google_trans_client"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/httpserver"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Cfg) {
	l := logger.New(cfg.Level)

	// Adapters

	// Google translate adapter
	client, err := googletransclient.New(cfg.DictionaryAPIURL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - googletransclient.New: %w", err))
	}
	g := googletrans.New(client, cfg.DefaultSrcLang, cfg.DefaultTrgtLang)

	// Db adapter
	pool, err := postgres.New(cfg.PGURL, cfg.MaxPoolSize)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	db := wordpostgres.New(pool)

	// Usecase/business logic
	service := usecase.New(db, g)

	// Register handlers
	c := chi.NewRouter()
	v1.Register(c, service, l, cfg.RateLimit, cfg.RateWindow, cfg.AllowedOrigins)

	// Start server
	serverCtx, cancelServerCtx := context.WithCancel(context.Background())
	defer cancelServerCtx()

	server := httpserver.New(cfg.Host+cfg.Port, cfg.WriteTimeout, cfg.ReadTimeout, c)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		server.Start()
	}()

	go func() {
		<-interrupt
		shutDownCtx, cancelShutDownCtx := context.WithTimeout(serverCtx, cfg.ShutdownTimeout)
		defer cancelShutDownCtx()
		if err := server.ShutDown(shutDownCtx); err != nil {
			l.Fatal(fmt.Errorf("app - Run - server.ShutDown: %w", err))
		}
		l.Info("Server was gracefully shutdown")
		cancelServerCtx()
	}()

	l.Info(fmt.Sprintf("Server started at %s%s", cfg.Host, cfg.Port))
	<-serverCtx.Done()
}
