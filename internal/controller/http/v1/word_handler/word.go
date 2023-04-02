// Package word implements port layer.
package wordhadnler

import (
	"context"
	"net/http"
	"time"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-playground/validator/v10"
)

type (
	WordService interface {
		AddWord(ctx context.Context, collection entity.Collection) error
		DeleteWord(ctx context.Context, collection entity.Collection) error
		UserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error)
		UpdateLearnInterval(ctx context.Context, collection entity.Collection) error
	}

	wordHandler struct {
		wordService WordService
		logger      logger.Interface
		v           *validator.Validate
	}
)

// Register func registers routes for chi.Mux router.
func Register(c *chi.Mux, srv WordService, l logger.Interface, rateLimit int,
	rateWindow time.Duration, allowedOrigings []string, allowedHeaders []string,
	defaultCorsDuration uint) {
	h := &wordHandler{
		wordService: srv,
		logger:      l,
		v:           validator.New(),
	}

	c.Use(middleware.RequestID)
	c.Use(middleware.Logger)
	c.Use(middleware.Recoverer)
	c.Use(httprate.Limit(
		rateLimit,
		rateWindow,
		httprate.WithLimitHandler(h.rateLimitResponse),
	))
	c.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: allowedOrigings,
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodDelete,
				http.MethodPost,
				http.MethodPut,
			},
			AllowedHeaders:   allowedHeaders,
			AllowCredentials: true,
			MaxAge:           int(defaultCorsDuration),
		},
	))
	c.Use(middleware.SetHeader("Content-Type", "application/json"))
	c.Use(h.jwtAuthenticator)
	c.Route("/v1", func(r chi.Router) {
		r.Route("/words", func(r chi.Router) {
			r.Delete("/", h.deleteWord)
			r.Put("/", h.updateLearnInterval)
			r.Get("/", h.userWords)
			r.Post("/", h.addWord)
		})
	})
}
