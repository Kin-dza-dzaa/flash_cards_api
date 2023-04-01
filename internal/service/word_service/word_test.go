package wordservice

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	wordrepomock "github.com/Kin-dza-dzaa/flash_cards_api/internal/service/word_service/word_repo_mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WordService_Suite struct {
	suite.Suite
	ws           *WordService
	wordRepoMock *wordrepomock.WordRepository
	trMock       *wordrepomock.Translator
	ctx          context.Context
}

func (s *WordService_Suite) SetupSuite() {
	s.ctx = context.Background()

	db := wordrepomock.NewWordPostgres(s.T())
	tr := wordrepomock.NewTranslator(s.T())

	s.wordRepoMock = db
	s.trMock = tr

	s.ws = New(db, tr)
}

func (s *WordService_Suite) Test_AddWord() {
	testCases := []struct {
		Name      string
		Coll      entity.Collection
		SetUpMock func()
		ExpectErr bool
	}{
		{
			Name: "Add new word",
			Coll: entity.Collection{
				Word:   "some_word",
				UserID: "12345",
				Name:   "some_coll",
			},
			SetUpMock: func() {
				s.wordRepoMock.On("IsWordInCollection", s.ctx, mock.Anything).Once().Return(false, nil)
				s.wordRepoMock.On("IsTransInDB", s.ctx, mock.Anything).Once().Return(false, nil)
				s.trMock.On("Translate", mock.Anything).Once().
					Return(entity.WordTrans{}, nil)
				s.wordRepoMock.On("AddTranslation", s.ctx, mock.Anything).Once().
					Return(nil)
				s.wordRepoMock.On("AddWord", s.ctx, mock.Anything).Once().Return(nil)
			},
			ExpectErr: false,
		},
		{
			Name: "Add existing word",
			Coll: entity.Collection{
				Word:   "some_word",
				UserID: "12345",
				Name:   "some_coll",
			},
			SetUpMock: func() {
				s.wordRepoMock.On("IsWordInCollection", s.ctx, mock.Anything).Once().Return(true, nil)
			},
			ExpectErr: false,
		},
		{
			Name: "Add new word but in DB",
			Coll: entity.Collection{
				Word:   "some_word",
				UserID: "12345",
				Name:   "some_coll",
			},
			SetUpMock: func() {
				s.wordRepoMock.On("IsWordInCollection", s.ctx, mock.Anything).Once().Return(false, nil)
				s.wordRepoMock.On("IsTransInDB", s.ctx, mock.Anything).Once().Return(true, nil)
				s.wordRepoMock.On("AddWord", s.ctx, mock.Anything).Once().Return(nil)
			},
			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			err := s.ws.AddWord(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func (s *WordService_Suite) Test_UserWords() {
	testCases := []struct {
		Name      string
		Coll      entity.Collection
		SetUpMock func()
		ExpectErr bool
		Res       *entity.UserWords
	}{
		{
			Name: "Add new word",
			Coll: entity.Collection{
				Word:   "some_word",
				UserID: "12345",
				Name:   "some_coll",
			},
			SetUpMock: func() {
				s.wordRepoMock.On("UserWords", s.ctx, mock.Anything).Once().
					Return(new(entity.UserWords), nil)
			},
			ExpectErr: false,
			Res:       new(entity.UserWords),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			userWords, err := s.ws.UserWords(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
			s.Assert().Equal(tc.Res, userWords, "User words must be equal")
		})
	}
}

func (s *WordService_Suite) Test_UpdateLearnInterval() {
	testCases := []struct {
		Name      string
		Coll      entity.Collection
		SetUpMock func()
		ExpectErr bool
	}{
		{
			Name: "Add new word",
			Coll: entity.Collection{
				Word:   "some_word",
				UserID: "12345",
				Name:   "some_coll",
			},
			SetUpMock: func() {
				s.wordRepoMock.On("UpdateLearnInterval", s.ctx, mock.Anything).Once().
					Return(nil)
			},
			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			err := s.ws.UpdateLearnInterval(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func (s *WordService_Suite) Test_DeleteWord() {
	testCases := []struct {
		Name      string
		Coll      entity.Collection
		SetUpMock func()
		ExpectErr bool
	}{
		{
			Name: "Add new word",
			Coll: entity.Collection{
				Word:   "some_word",
				UserID: "12345",
				Name:   "some_coll",
			},
			SetUpMock: func() {
				s.wordRepoMock.On("DeleteWord", s.ctx, mock.Anything).Once().
					Return(nil)
			},
			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			err := s.ws.DeleteWord(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func Test_StartUseCase_Suite(t *testing.T) {
	suite.Run(t, new(WordService_Suite))
}
