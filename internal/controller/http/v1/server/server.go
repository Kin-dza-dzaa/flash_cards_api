// Package server implements HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	"golang.org/x/exp/slog"
)

type Server struct {
	cfg config.Cfg
	l   *slog.Logger
	*http.Server
}

func (srv *Server) Start(ctx context.Context) <-chan struct{} {
	serverCtx, cancelServerCtx := context.WithCancel(ctx)

	go func() {
		defer cancelServerCtx()
		err := srv.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			srv.l.Error(
				"Unable to start server",
				slog.String("error", err.Error()),
			)
		}
	}()

	srv.gracefullyShutdown(serverCtx)
	srv.l.Info(fmt.Sprintf("Server started at %s", srv.cfg.HTTP.Addr))

	return serverCtx.Done()
}

func (srv *Server) gracefullyShutdown(ctx context.Context) {
	go func(ctx context.Context) {
		interrupt := make(chan os.Signal, 1)
		defer close(interrupt)
		signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case <-interrupt:
			shutdownCtx, cancelShutdownCtx := context.WithTimeout(ctx, srv.cfg.HTTP.ShutdownTimeout)
			defer cancelShutdownCtx()
			err := srv.Server.Shutdown(shutdownCtx)
			if err != nil {
				srv.l.Error(
					"Error while shutting down server",
					slog.String("error", err.Error()),
				)
				return
			}
			srv.l.Info("Server was gracefully shutdown")
			return
		case <-ctx.Done():
			shutdownCtx, cancelShutdownCtx := context.WithTimeout(ctx, srv.cfg.HTTP.ShutdownTimeout)
			defer cancelShutdownCtx()
			err := srv.Server.Shutdown(shutdownCtx)
			if err != nil {
				srv.l.Error(
					"Error while shutting down server",
					slog.String("error", err.Error()),
				)
				return
			}
			srv.l.Info("Server was gracefully shutdown")
			return
		}
	}(ctx)
}

func New(cfg config.Cfg, l *slog.Logger, h http.Handler) *Server {
	srv := &http.Server{
		Addr:         cfg.HTTP.Addr,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		Handler:      h,
	}

	return &Server{
		Server: srv,
		cfg:    cfg,
		l:      l,
	}
}
