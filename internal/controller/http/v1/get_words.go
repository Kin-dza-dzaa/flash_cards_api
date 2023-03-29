package v1

import (
	"fmt"
	"net/http"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

func (h *wordHandler) getWords(w http.ResponseWriter, r *http.Request) {
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
	collection := &entity.Collection{
		UserID: userID,
	}

	words, err := h.srv.GetUserWords(r.Context(), collection)
	if err != nil {
		h.logger.Error(fmt.Errorf("wordHandler - getWords - h.service.GetUserWords: %w", err))
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
			Message:   "success",
			UserWords: words,
		},
	)
}
