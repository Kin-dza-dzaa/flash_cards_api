package wordhadnler

import (
	"testing"

	wordservicemock "github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1/word_handler/word_service_mock"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

// Base suite for wordHandler
type wordHandler_Suite struct {
	suite.Suite
	srv *wordservicemock.WordService
	h   *wordHandler
}

func (s *wordHandler_Suite) SetupSuite() {
	s.srv = wordservicemock.NewWordService(s.T())
	s.h = &wordHandler{
		wordService: s.srv,
		logger:      logger.New("debug"),
		v:           validator.New(),
	}
}

func Test_wordHandler_Suite(t *testing.T) {
	suite.Run(t, new(wordHandler_Suite))
}
