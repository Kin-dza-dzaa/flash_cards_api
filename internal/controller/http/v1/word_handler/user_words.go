package wordhadnler

import (
	"fmt"
	"net/http"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

func (h *wordHandler) userWords(w http.ResponseWriter, r *http.Request) {
	userID, err := fromCtx[string](r.Context(), "user_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.encodeResponse(w,
			httpResponse{
				Status:  http.StatusUnauthorized,
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
		h.logger.Error(fmt.Errorf("wordHandler - userWords - h.service.UserWords: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeResponse(w,
			httpResponse{
				Path:    r.URL.Path,
				Status:  http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}

	h.encodeResponse(w,
		httpResponse{
			Path:      r.URL.Path,
			Status:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
			UserWords: words,
		},
	)
}
