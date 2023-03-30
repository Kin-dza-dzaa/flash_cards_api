package usecase

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/usecase/adaptermock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WordUseCaseTestSuite struct {
	suite.Suite
	uc     *WordService
	dbMock *adaptermock.Database
	trMock *adaptermock.Tranlsator
	ctx    context.Context
}

func (s *WordUseCaseTestSuite) SetupSuite() {
	s.ctx = context.Background()

	db := adaptermock.NewDatabase(s.T())
	tr := adaptermock.NewTranlsator(s.T())

	s.dbMock = db
	s.trMock = tr

	s.uc = New(db, tr)
}

func (s *WordUseCaseTestSuite) TestAddWord() {
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
				s.dbMock.On("IsWordInCollection", s.ctx, mock.Anything).Once().Return(false, nil)
				s.dbMock.On("IsTransInDB", s.ctx, mock.Anything).Once().Return(false, nil)
				s.trMock.On("Translate", mock.Anything).Once().
					Return(entity.WordTrans{}, nil)
				s.dbMock.On("AddTranslation", s.ctx, mock.Anything).Once().
					Return(nil)
				s.dbMock.On("AddWordToCollection", s.ctx, mock.Anything).Once().Return(nil)
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
				s.dbMock.On("IsWordInCollection", s.ctx, mock.Anything).Once().Return(true, nil)
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
				s.dbMock.On("IsWordInCollection", s.ctx, mock.Anything).Once().Return(false, nil)
				s.dbMock.On("IsTransInDB", s.ctx, mock.Anything).Once().Return(true, nil)
				s.dbMock.On("AddWordToCollection", s.ctx, mock.Anything).Once().Return(nil)
			},
			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			err := s.uc.AddWord(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func (s *WordUseCaseTestSuite) TestGetUserWords() {
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
				s.dbMock.On("GetUserWords", s.ctx, mock.Anything).Once().
					Return(new(entity.UserWords), nil)
			},
			ExpectErr: false,
			Res:       new(entity.UserWords),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			userWords, err := s.uc.GetUserWords(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
			s.Assert().Equal(tc.Res, userWords, "User words must be equal")
		})
	}
}

func (s *WordUseCaseTestSuite) TestUpdateLearnInterval() {
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
				s.dbMock.On("UpdateLearnInterval", s.ctx, mock.Anything).Once().
					Return(nil)
			},
			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			err := s.uc.UpdateLearnInterval(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func (s *WordUseCaseTestSuite) TestDeleteWordFromCollection() {
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
				s.dbMock.On("DeleteWordFromCollection", s.ctx, mock.Anything).Once().
					Return(nil)
			},
			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock()
			err := s.uc.DeleteWordFromCollection(s.ctx, tc.Coll)
			if tc.ExpectErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func TestStartUseCaseSuite(t *testing.T) {
	suite.Run(t, new(WordUseCaseTestSuite))
}
