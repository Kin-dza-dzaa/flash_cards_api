// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	wordhadnler "github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1/word_handler"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/repository"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/service"
	googletransclient "github.com/Kin-dza-dzaa/flash_cards_api/pkg/google_trans_client"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func Run(cfg config.Cfg) {
	l := logger.New(cfg.Logger.Level)

	// Adapters

	// Google translate HTTP 2.0 client
	client, err := googletransclient.New(cfg.GoogleApi.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - googletransclient.New: %w", err))
	}
	// Postgres pool
	pool, err := postgres.New(cfg.PG.URL, cfg.PG.MaxPoolSize)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pool.Close()

	// Repo
	r := repository.New(client, pool, cfg.GoogleApi.DefaultSrcLang, cfg.GoogleApi.DefaultTrgtLang)

	// Usecase/business logic
	s := service.New(r, r)

	// Register routes
	c := chi.NewRouter()
	wordhadnler.Register(c, s, l, cfg.HTTP.RateLimit, cfg.HTTP.RateWindow,
		cfg.HTTP.AllowedOrigins, cfg.HTTP.AllowedHeaders, cfg.HTTP.DefaultCorsDuration)

	// Configure server
	srv := http.Server{
		Addr:         cfg.HTTP.Addr,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		Handler:      c,
	}

	// Server start-up
	serverCtx, cancelServerCtx := context.WithCancel(context.Background())

	go func() {
		defer cancelServerCtx()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error(fmt.Errorf("App - Run - srv.ListenAndServer: %w", err))
			return
		}
		l.Info("Server was gracefully shutdown")
	}()

	// Gracefull shutdown
	interrupt := make(chan os.Signal, 1)
	defer close(interrupt)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		select {
		case <-interrupt:
			shutdownCtx, cancelShutDownCtx := context.WithTimeout(serverCtx, cfg.HTTP.ShutdownTimeout)
			defer cancelShutDownCtx()
			if err := srv.Shutdown(shutdownCtx); err != nil {
				l.Error(fmt.Errorf("app - Run - server.ShutDown: %w", err))
			}
		case <-serverCtx.Done():
		}
	}()

	l.Info(fmt.Sprintf("Server started at %s", cfg.HTTP.Addr))
	<-serverCtx.Done()
}
