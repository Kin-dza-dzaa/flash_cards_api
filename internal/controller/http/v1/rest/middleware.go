package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"golang.org/x/exp/slog"
	"google.golang.org/api/idtoken"
)

func (h *WordHandler) jwtAuthenticator(next http.Handler) http.Handler {
	respond := func(w http.ResponseWriter, r *http.Request) {
		h.encode(w,
			http.StatusUnauthorized,
			httpResponse{
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusUnauthorized),
			})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header.
		token := jwtauth.TokenFromHeader(r)
		if token == "" {
			respond(w, r)
			return
		}

		// Validate token, makes HTTP request to get PK from google.
		p, err := idtoken.Validate(r.Context(), token, "")
		if err != nil {
			respond(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(inCtx(r.Context(), "user_id", p.Subject)))
	})
}

func (h *WordHandler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		defer func() {
			path := r.URL.Path
			h.logger.InfoCtx(
				r.Context(),
				"request",
				slog.Int64("elapsed", time.Since(t).Nanoseconds()),
				slog.String("path", path),
				slog.String("method", r.Method),
			)
		}()

		next.ServeHTTP(w, r)
	})
}
