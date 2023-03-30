// Package v1 implements port layer.
package v1

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

const (
	defualtCorsDuration = 5
)

type (
	service interface {
		AddWord(ctx context.Context, collection entity.Collection) error
		DeleteWordFromCollection(ctx context.Context, collection entity.Collection) error
		GetUserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error)
		UpdateLearnInterval(ctx context.Context, collection entity.Collection) error
	}

	wordHandler struct {
		srv    service
		logger logger.Interface
		v      *validator.Validate
	}
)

// Register func registers routes for chi.Mux router.
func Register(c *chi.Mux, srv service, l logger.Interface, rateLimit int,
	rateWindow time.Duration, allowedOrigings []string) {
	h := &wordHandler{
		srv:    srv,
		logger: l,
		v:      validator.New(),
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
			AllowedHeaders: []string{
				"Content-Type",
				"Authorization",
			},
			AllowCredentials: true,
			MaxAge:           defualtCorsDuration,
		},
	))
	c.Use(middleware.SetHeader("Content-Type", "application/json"))
	c.Use(h.jwtAuthenticator)
	c.Route("/v1", func(r chi.Router) {
		r.Route("/words", func(r chi.Router) {
			r.Delete("/", h.deleteWordFromCollection)
			r.Put("/", h.updateLearnInterval)
			r.Get("/", h.getWords)
			r.Post("/", h.addWord)
		})
	})
}
