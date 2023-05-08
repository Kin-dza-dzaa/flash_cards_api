// Package rest implements port layer.
package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-playground/validator/v10"
	"github.com/riandyrn/otelchi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/exp/slog"
)

const (
	otelName     = "github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1/rest"
	userIDCtxKey = "user_id"
)

type (
	wordService interface {
		AddWord(ctx context.Context, collection entity.Collection) error
		DeleteWord(ctx context.Context, collection entity.Collection) error
		UserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error)
		UpdateLearnInterval(ctx context.Context, collection entity.Collection) error
	}
)

type WordHandler struct {
	wordService wordService
	logger      *slog.Logger
	v           *validator.Validate
}

type UpdateLearnIntervalRequest struct {
	Word           string        `json:"word" validate:"required"`
	CollectionName string        `json:"collection_name" validate:"required"`
	LastRepeat     time.Time     `json:"last_repeat" validate:"required"`
	TimeDiff       time.Duration `json:"time_diff" validate:"required"`
}

type AddWordRequest struct {
	Word           string        `json:"word" validate:"required"`
	CollectionName string        `json:"collection_name" validate:"required"`
	LastRepeat     time.Time     `json:"last_repeat" validate:"required"`
	TimeDiff       time.Duration `json:"time_diff"`
}

type DeleteWordRequest struct {
	Word           string `json:"word" validate:"required"`
	CollectionName string `json:"collection_name" validate:"required"`
}

//	@title			Flash cards API
//	@version		0.3.4
//	@description	REST API for word and collections of a user.

//	@contact.name	API Support

// @host		localhost:8000
// @BasePath	/v1
func (h *WordHandler) Register(c *chi.Mux, cfg config.Cfg) {
	// Engine.
	c.Use(middleware.RequestID)
	c.Use(h.logRequest)
	c.Use(middleware.Recoverer)
	c.Use(httprate.Limit(
		cfg.HTTP.RateLimit,
		cfg.HTTP.RateWindow,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}),
	))
	c.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   cfg.HTTP.AllowedOrigins,
			AllowedMethods:   cfg.HTTP.AllowedMethods,
			AllowedHeaders:   cfg.HTTP.AllowedHeaders,
			AllowCredentials: cfg.HTTP.AllowCredentials,
			MaxAge:           int(cfg.HTTP.DefaultCorsDuration),
		},
	))

	// Routes.
	c.Get("/swagger/*", httpSwagger.Handler())

	c.Route("/v1", func(r chi.Router) {
		r.Use(h.jwtAuthenticator)
		r.Use(otelchi.Middleware("flash-cards-api-server"))
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Route("/words", func(r chi.Router) {
			r.Delete("/", h.deleteWord)
			r.Put("/", h.updateLearnInterval)
			r.Get("/", h.userWords)
			r.Post("/", h.addWord)
		})
	})
}

// List user words
//
//	@Summary		Get user words.
//	@Description	Gets user words that put together in collections.
//	@Tags			words
//	@Produce		json
//	@Success		200	{object}	entity.UserWords	"User words"
//	@Failure		401	{object}	httpResponse		"Unauthorized"
//	@Failure		500	{object}	httpResponse		"Internal error"
//	@Router			/words [get]
func (h *WordHandler) userWords(w http.ResponseWriter, r *http.Request) {
	userID := fromCtx(r.Context(), userIDCtxKey)
	if userID == "" {
		h.encode(
			w,
			http.StatusUnauthorized,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusUnauthorized),
			})
		return
	}
	collection := entity.Collection{
		UserID: userID,
	}

	words, err := h.wordService.UserWords(r.Context(), collection)
	if err != nil {
		h.logger.ErrorCtx(
			r.Context(),
			"Internal error",
			slog.String("error", fmt.Errorf("wordHandler - userWords - h.service.UserWords: %w", err).Error()),
		)
		h.encode(
			w,
			http.StatusInternalServerError,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusInternalServerError),
			},
		)

		_, span := otel.Tracer(otelName).Start(r.Context(), "WordHandler - userWords - Error")
		defer span.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	h.encode(
		w,
		http.StatusOK,
		words,
	)
}

// Update learn internal of a word.
//
//	@Summary	Updates learn interval for a given word.
//	@Tags		words
//	@Accept		json
//	@Produce	json
//	@Param		WordInfo	body		UpdateLearnIntervalRequest	true	"Word, collection name with learn intervals"
//	@Success	200			{object}	httpResponse				"Interval was updated"
//	@Failure	400			{object}	httpResponse				"Wrong JSON format"
//	@Failure	401			{object}	httpResponse				"Unauthorized"
//	@Failure	500			{object}	httpResponse				"Internal error"
//	@Router		/words [put]
func (h *WordHandler) updateLearnInterval(w http.ResponseWriter, r *http.Request) {
	userID := fromCtx(r.Context(), userIDCtxKey)
	if userID == "" {
		h.encode(
			w,
			http.StatusUnauthorized,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusUnauthorized),
			})
		return
	}

	var req UpdateLearnIntervalRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Path:    r.URL.Path,
				Message: wrongJSONFormat,
			},
		)
		return
	}

	if err := h.v.Struct(req); err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusBadRequest),
			})
		return
	}

	err = h.wordService.UpdateLearnInterval(
		r.Context(),
		entity.Collection{
			UserID:     userID,
			Name:       req.CollectionName,
			Word:       req.Word,
			LastRepeat: req.LastRepeat,
			TimeDiff:   req.TimeDiff,
		},
	)
	if err != nil {
		h.encode(
			w,
			http.StatusInternalServerError,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusInternalServerError),
			},
		)

		_, span := otel.Tracer(otelName).Start(r.Context(), "WordHandler - updateLearnInterval - Error")
		defer span.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	h.encode(
		w,
		http.StatusOK,
		httpResponse{
			Path:    r.URL.Path,
			Message: http.StatusText(http.StatusOK),
		})
}

// Delete word from collection.
//
//	@Summary	Deletes given word from a collection.
//	@Tags		words
//	@Accept		json
//	@Produce	json
//	@Param		WordInfo	body		DeleteWordRequest	true	"Word and collection name"
//	@Success	200			{object}	httpResponse		"Word was deleted"
//	@Failure	400			{object}	httpResponse		"Wrong JSON format"
//	@Failure	401			{object}	httpResponse		"Unauthorized"
//	@Failure	500			{object}	httpResponse		"Internal error"
//	@Router		/words [delete]
func (h *WordHandler) deleteWord(w http.ResponseWriter, r *http.Request) {
	userID := fromCtx(r.Context(), userIDCtxKey)
	if userID == "" {
		h.encode(
			w,
			http.StatusUnauthorized,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusUnauthorized),
			})
		return
	}

	var req DeleteWordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Path:    r.URL.Path,
				Message: wrongJSONFormat,
			},
		)
		return
	}

	if err := h.v.Struct(req); err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusBadRequest),
			})
		return
	}

	err = h.wordService.DeleteWord(
		r.Context(),
		entity.Collection{
			UserID: userID,
			Name:   req.CollectionName,
			Word:   req.Word,
		},
	)
	if err != nil {
		h.logger.ErrorCtx(
			r.Context(),
			"Internal error",
			slog.String("error", fmt.Errorf("wordHandler - deleteWord - h.service.deleteWord: %w", err).Error()),
		)
		h.encode(
			w,
			http.StatusInternalServerError,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusInternalServerError),
			},
		)

		_, span := otel.Tracer(otelName).Start(r.Context(), "WordHandler - deleteWord - Error")
		defer span.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	h.encode(
		w,
		http.StatusOK,
		httpResponse{
			Path:    r.URL.Path,
			Message: http.StatusText(http.StatusOK),
		})
}

// Add word to collection.
//
//	@Summary	Adds a word to a given collection.
//	@Tags		words
//	@Accept		json
//	@Produce	json
//	@Param		WordInfo	body		AddWordRequest	true	"Word, collection name with learn intervals"
//	@Success	201			{object}	httpResponse	"Word was added to collection"
//	@Failure	400			{object}	httpResponse	"Wrong JSON format"
//	@Failure	401			{object}	httpResponse	"Unauthorized"
//	@Failure	403			{object}	httpResponse	"Word not supported"
//	@Failure	500			{object}	httpResponse	"Internal error"
//	@Router		/words [post]
func (h *WordHandler) addWord(w http.ResponseWriter, r *http.Request) {
	userID := fromCtx(r.Context(), userIDCtxKey)
	if userID == "" {
		h.encode(
			w,
			http.StatusUnauthorized,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusUnauthorized),
			})
		return
	}

	var req AddWordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Path:    r.URL.Path,
				Message: wrongJSONFormat,
			},
		)
		return
	}

	if err := h.v.Struct(req); err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusBadRequest),
			})
		return
	}

	err = h.wordService.AddWord(
		r.Context(),
		entity.Collection{
			UserID:     userID,
			Name:       req.CollectionName,
			Word:       req.Word,
			LastRepeat: req.LastRepeat,
			TimeDiff:   req.TimeDiff,
		},
	)
	if err != nil {
		if errors.Is(err, entity.ErrWordNotSupported) {
			h.encode(
				w,
				http.StatusForbidden,
				httpResponse{
					Path:    r.URL.Path,
					Message: entity.ErrWordNotSupported.Error(),
				},
			)
			return
		}

		h.logger.ErrorCtx(
			r.Context(),
			"Internal error",
			slog.String("error", fmt.Errorf("wordHandler - getWords - h.service.AddWord: %w", err).Error()),
		)
		h.encode(
			w,
			http.StatusInternalServerError,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusInternalServerError),
			},
		)

		_, span := otel.Tracer(otelName).Start(r.Context(), "WordHandler - addWord - Error")
		defer span.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	h.encode(
		w,
		http.StatusCreated,
		httpResponse{
			Path:    r.URL.Path,
			Message: http.StatusText(http.StatusCreated),
		})
}

func NewWordHandler(wordService wordService, l *slog.Logger) *WordHandler {
	h := &WordHandler{
		wordService: wordService,
		logger:      l,
		v:           validator.New(),
	}

	return h
}
