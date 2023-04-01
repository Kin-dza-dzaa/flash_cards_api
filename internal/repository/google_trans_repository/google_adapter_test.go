package googletransrepository

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	googletransclient "github.com/Kin-dza-dzaa/flash_cards_api/pkg/google_trans_client"
	"github.com/stretchr/testify/suite"
)

type GoogleAdapterTestSuite struct {
	suite.Suite
	tr *GoogleTranslate
}

func (s *GoogleAdapterTestSuite) SetupSuite() {
	cfg, err := config.ReadConfig()
	if err != nil {
		s.FailNow(err.Error())
	}

	gc, err := googletransclient.New(cfg.GoogleApi.URL)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.tr = New(gc, cfg.GoogleApi.DefaultSrcLang, cfg.GoogleApi.DefaultTrgtLang)
}

// Test makes real calls to google translate api
func (s *GoogleAdapterTestSuite) TestTranslate() {
	testCases := []struct {
		Name string
		Word string
		Err  error
	}{
		{
			Name: "Check unsupported word",
			Word: "bad_word!!!!!@#!@$#!%#",
			Err:  entity.ErrWordNotSupported,
		},
		{
			Name: "Check valid word",
			Word: "lead",
			Err:  nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			wordTrans, err := s.tr.Translate(tc.Word)
			s.Assert().Equal(tc.Err, err, "Errors must be equal")
			if tc.Err == nil {
				s.Assert().NotEmpty(wordTrans, "Word translation must be not empty")
			}
		})
	}
}

func TestStartGoogleAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(GoogleAdapterTestSuite))
}
