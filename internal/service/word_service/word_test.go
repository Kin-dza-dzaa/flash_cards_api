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
}

func (s *WordService_Suite) SetupSuite() {
	db := wordrepomock.NewWordPostgres(s.T())
	tr := wordrepomock.NewTranslator(s.T())

	s.wordRepoMock = db
	s.trMock = tr

	s.ws = New(db, tr)
}

func (s *WordService_Suite) Test_AddWord() {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		Coll entity.Collection
	}

	testCases := []struct {
		Name      string
		Args      args
		SetUpMock func(args args)
		WantErr   bool
	}{
		{
			Name: "Add new word",
			Args: args{
				ctx: ctx,
				Coll: entity.Collection{
					Name: "some_name",
					Word: "Some_words",
				},
			},
			SetUpMock: func(args args) {
				s.wordRepoMock.On("IsWordInCollection", args.ctx, args.Coll).Once().Return(false, nil)
				s.wordRepoMock.On("IsTransInDB", args.ctx, args.Coll).Once().Return(false, nil)
				s.trMock.On("Translate", args.Coll.Word).Once().
					Return(entity.WordTrans{}, nil)
				s.wordRepoMock.On("AddTranslation", args.ctx, mock.Anything).Once().
					Return(nil)
				s.wordRepoMock.On("AddWord", args.ctx, args.Coll).Once().Return(nil)
			},
			WantErr: false,
		},
		{
			Name: "Add existing word",
			Args: args{
				ctx: ctx,
				Coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			SetUpMock: func(args args) {
				s.wordRepoMock.On("IsWordInCollection", args.ctx, args.Coll).Once().Return(true, nil)
			},
			WantErr: false,
		},
		{
			Name: "Add new word but in DB",
			Args: args{
				ctx: ctx,
				Coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			SetUpMock: func(args args) {
				s.wordRepoMock.On("IsWordInCollection", args.ctx, args.Coll).Once().Return(false, nil)
				s.wordRepoMock.On("IsTransInDB", args.ctx, args.Coll).Once().Return(true, nil)
				s.wordRepoMock.On("AddWord", args.ctx, args.Coll).Once().Return(nil)
			},
			WantErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock(tc.Args)
			err := s.ws.AddWord(tc.Args.ctx, tc.Args.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func (s *WordService_Suite) Test_UserWords() {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		Coll entity.Collection
	}
	testCases := []struct {
		Name      string
		Args      args
		SetUpMock func(args args)
		WantErr   bool
		Res       *entity.UserWords
	}{
		{
			Name: "Add new word",
			Args: args{
				ctx: ctx,
				Coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			SetUpMock: func(args args) {
				s.wordRepoMock.On("UserWords", args.ctx, args.Coll).Once().
					Return(new(entity.UserWords), nil)
			},
			WantErr: false,
			Res:     new(entity.UserWords),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock(tc.Args)
			userWords, err := s.ws.UserWords(tc.Args.ctx, tc.Args.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
			s.Assert().Equal(tc.Res, userWords, "User words must be equal")
		})
	}
}

func (s *WordService_Suite) Test_UpdateLearnInterval() {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		Coll entity.Collection
	}
	testCases := []struct {
		Name      string
		Args      args
		SetUpMock func(args args)
		WantErr   bool
	}{
		{
			Name: "Add new word",
			Args: args{
				ctx: ctx,
				Coll: entity.Collection{
					Word:   "some_word",
					UserID: "12345",
					Name:   "some_coll",
				},
			},
			SetUpMock: func(args args) {
				s.wordRepoMock.On("UpdateLearnInterval", args.ctx, args.Coll).Once().
					Return(nil)
			},
			WantErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock(tc.Args)
			err := s.ws.UpdateLearnInterval(tc.Args.ctx, tc.Args.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func (s *WordService_Suite) Test_DeleteWord() {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		Coll entity.Collection
	}
	testCases := []struct {
		Name      string
		Args      args
		SetUpMock func(args args)
		ExpectErr bool
	}{
		{
			Name: "Add new word",
			Args: args{
				ctx: ctx,
				Coll: entity.Collection{
					Word:   "some_word",
					UserID: "12345",
					Name:   "some_coll",
				},
			},
			SetUpMock: func(args args) {
				s.wordRepoMock.On("DeleteWord", args.ctx, args.Coll).Once().
					Return(nil)
			},
			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tc.SetUpMock(tc.Args)
			err := s.ws.DeleteWord(tc.Args.ctx, tc.Args.Coll)
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
