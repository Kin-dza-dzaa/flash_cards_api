package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

const (
	wrongJsonFormat = "wrong json format"
)

type httpResponse struct {
	Path      string            `json:"path"`
	Status    int               `json:"status"`
	Message   string            `json:"message"`
	UserWords *entity.UserWords `json:"user_words,omitempty"`
}

// Parses json, if error happens method will return it.
func (h *wordHandler) decodeCollection(r *http.Request) (*entity.Collection, error) {
	collection := new(entity.Collection)
	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// Encodes in w stream, before calling that function use w.WriteHeader method
// After calling that function you shouldn't write to w.
func (h *wordHandler) encodeResponse(w http.ResponseWriter, response httpResponse) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(fmt.Errorf("wordHandler - encodeResponse - Encode: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
