package googletrans

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	googletransclient "github.com/Kin-dza-dzaa/flash_cards_api/pkg/google_trans_client"
	"github.com/stretchr/testify/suite"
)

type GoogleAdapterTestSuite struct {
	suite.Suite
	tr *GoogleTranslate
}

func (s *GoogleAdapterTestSuite) SetupSuite() {
	const (
		defaultSrcLang  = "en"
		defaultTrgtLang = "ru"
		apiURL          = "https://translate.google.com/_/TranslateWebserverUi/data/batchexecute"
	)

	gc, err := googletransclient.New(apiURL)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.tr = New(gc, defaultSrcLang, defaultTrgtLang)
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
