package wordhadnler

import (
	"fmt"
	"net/http"
)

func (h *wordHandler) deleteWord(w http.ResponseWriter, r *http.Request) {
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

	err = h.wordService.DeleteWord(r.Context(), collection)
	if err != nil {
		h.logger.Error(fmt.Errorf("wordHandler - deleteWord - h.service.deleteWord: %w", err))
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
