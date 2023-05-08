package logger

import (
	"context"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

// CustomJSONHandler tries to get requestID from ctx and use it in structured log.
type CustomJSONHandler struct {
	*slog.JSONHandler
}

func (h *CustomJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	rID := middleware.GetReqID(ctx)
	if rID != "" {
		r.AddAttrs(slog.String("req_id", rID))
	}

	return h.JSONHandler.Handle(ctx, r)
}

func New(level slog.Leveler) *slog.Logger {
	basicJSONHandler := slog.HandlerOptions{
		Level: level,
	}.NewJSONHandler(os.Stdout)

	logger := slog.New(
		&CustomJSONHandler{
			basicJSONHandler,
		},
	)
	return logger
}
