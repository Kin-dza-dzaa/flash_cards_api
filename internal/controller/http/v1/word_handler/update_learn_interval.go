package wordhadnler

import (
	"net/http"
)

func (h *wordHandler) updateLearnInterval(w http.ResponseWriter, r *http.Request) {
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
				Message: wrongJsonFormat,
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

	if err := h.wordService.UpdateLearnInterval(r.Context(), collection); err != nil {
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
