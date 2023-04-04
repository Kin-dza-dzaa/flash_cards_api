package wordhadnler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

func (h *wordHandler) addWord(w http.ResponseWriter, r *http.Request) {
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

	collection, err := h.decodeCollection(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.encodeResponse(w,
			httpResponse{
				Path:    r.URL.Path,
				Status:  http.StatusBadRequest,
				Message: wrongJSONFormat,
			},
		)
		return
	}
	collection.UserID = userID

	if err := h.v.Struct(collection); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.encodeResponse(w,
			httpResponse{
				Status:  http.StatusBadRequest,
				Path:    r.URL.Path,
				Message: http.StatusText(http.StatusBadRequest),
			})
		return
	}

	err = h.wordService.AddWord(r.Context(), collection)
	if err != nil {
		if errors.Is(err, entity.ErrWordNotSupported) {
			w.WriteHeader(http.StatusBadRequest)
			h.encodeResponse(w,
				httpResponse{
					Path:    r.URL.Path,
					Status:  http.StatusBadRequest,
					Message: entity.ErrWordNotSupported.Error(),
				},
			)
			return
		}

		h.logger.Error(fmt.Errorf("wordHandler - getWords - h.service.AddWord: %w", err))
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
			Path:    r.URL.Path,
			Status:  http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		})
}
