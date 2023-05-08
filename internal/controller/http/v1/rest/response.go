package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
)

const (
	wrongJSONFormat = "wrong json format"
)

type httpResponse struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

// Encodes in w stream.
// After calling that function you shouldn't write to w.
func (h *WordHandler) encode(w http.ResponseWriter, status int, response interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(
			"JSON encoding error",
			slog.String("error", fmt.Errorf("wordHandler - encodeResponse - Encode: %w", err).Error()),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
