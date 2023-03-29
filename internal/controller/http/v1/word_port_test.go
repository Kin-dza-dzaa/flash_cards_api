package v1

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/srvmock"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

// Base suite for wordHandler
type wordHandler_Suite struct {
	suite.Suite
	srv *srvmock.Service
	h   *wordHandler
}

func (s *wordHandler_Suite) SetupSuite() {
	s.srv = srvmock.NewService(s.T())
	s.h = &wordHandler{
		srv:    s.srv,
		logger: logger.New("debug"),
		v:      validator.New(),
	}
}

func Test_wordHandler_Suite(t *testing.T) {
	suite.Run(t, new(wordHandler_Suite))
}
