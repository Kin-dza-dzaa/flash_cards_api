package wordhadnler

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"google.golang.org/api/idtoken"
)

func (h *wordHandler) jwtAuthenticator(next http.Handler) http.Handler {

	respond := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		h.encodeResponse(w,
			httpResponse{
				Status:  http.StatusUnauthorized,
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusUnauthorized),
			})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from Authorization header
		token := jwtauth.TokenFromHeader(r)
		if token == "" {
			respond(w, r)
			return
		}

		// validate token
		p, err := idtoken.Validate(r.Context(), token, "")
		if err != nil {
			respond(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(inCtx(r.Context(), "user_id", p.Subject)))
	})
}

func (h *wordHandler) rateLimitResponse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
}
