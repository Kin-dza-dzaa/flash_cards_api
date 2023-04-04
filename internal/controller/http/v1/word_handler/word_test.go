package wordhadnler

import (
	"testing"

	wordservicemock "github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1/word_handler/word_service_mock"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/go-playground/validator/v10"
)

func setupWordHandler(t *testing.T) (*wordHandler, *wordservicemock.WordService) {
	t.Helper()
	srvMock := wordservicemock.NewWordService(t)
	h := &wordHandler{
		wordService: srvMock,
		logger:      logger.New("debug"),
		v:           validator.New(),
	}
	return h, srvMock
}
