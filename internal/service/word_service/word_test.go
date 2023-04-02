package wordservice

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	wordrepomock "github.com/Kin-dza-dzaa/flash_cards_api/internal/service/word_service/word_repo_mock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
)

func setupWordService(t *testing.T) (*wordrepomock.WordRepository, *wordrepomock.Translator) {
	db := wordrepomock.NewWordPostgres(t)
	tr := wordrepomock.NewTranslator(t)

	return db, tr
}

func Test_AddWord(t *testing.T) {
	ctx := context.Background()
	wordRepoMock, trMock := setupWordService(t)
	wordService := New(wordRepoMock, trMock)

	type args struct {
		ctx  context.Context
		coll entity.Collection
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(args args)
		wantErr   bool
	}{
		{
			name: "Add new word",
			args: args{
				ctx: ctx,
				coll: entity.Collection{
					Name: "some_name",
					Word: "Some_words",
				},
			},
			setupMock: func(args args) {
				wordRepoMock.On("IsWordInCollection", args.ctx, args.coll).Once().Return(false, nil)
				wordRepoMock.On("IsTransInDB", args.ctx, args.coll).Once().Return(false, nil)
				trMock.On("Translate", args.coll.Word).Once().
					Return(entity.WordTrans{}, nil)
				wordRepoMock.On("AddTranslation", args.ctx, mock.Anything).Once().
					Return(nil)
				wordRepoMock.On("AddWord", args.ctx, args.coll).Once().Return(nil)
			},
		},
		{
			name: "Add existing word",
			args: args{
				ctx: ctx,
				coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			setupMock: func(args args) {
				wordRepoMock.On("IsWordInCollection", args.ctx, args.coll).Once().Return(true, nil)
			},
		},
		{
			name: "Add word that in DB",
			args: args{
				ctx: ctx,
				coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			setupMock: func(args args) {
				wordRepoMock.On("IsWordInCollection", args.ctx, args.coll).Once().Return(false, nil)
				wordRepoMock.On("IsTransInDB", args.ctx, args.coll).Once().Return(true, nil)
				wordRepoMock.On("AddWord", args.ctx, args.coll).Once().Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.args)
			err := wordService.AddWord(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
		})
	}
}

func Test_UserWords(t *testing.T) {
	ctx := context.Background()
	wordRepoMock, trMock := setupWordService(t)
	wordService := New(wordRepoMock, trMock)

	type args struct {
		ctx  context.Context
		coll entity.Collection
	}
	tests := []struct {
		name          string
		args          args
		setupMock     func(args args)
		wantErr       bool
		wantUserWords *entity.UserWords
	}{
		{
			name: "Add new word",
			args: args{
				ctx: ctx,
				coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			setupMock: func(args args) {
				wordRepoMock.On("UserWords", args.ctx, args.coll).Once().
					Return(new(entity.UserWords), nil)
			},
			wantErr:       false,
			wantUserWords: new(entity.UserWords),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.args)
			gotUserWords, err := wordService.UserWords(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
			if diff := cmp.Diff(gotUserWords, tt.wantUserWords); diff != "" {
				t.Fatalf("wanted: %v but got %v", tt.wantUserWords, gotUserWords)
			}
		})
	}
}

func Test_UpdateLearnInterval(t *testing.T) {
	ctx := context.Background()
	wordRepoMock, trMock := setupWordService(t)
	wordService := New(wordRepoMock, trMock)

	type args struct {
		ctx  context.Context
		coll entity.Collection
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(args args)
		wantErr   bool
	}{
		{
			name: "Add new word",
			args: args{
				ctx: ctx,
				coll: entity.Collection{
					Word:   "some_word",
					UserID: "12345",
					Name:   "some_coll",
				},
			},
			setupMock: func(args args) {
				wordRepoMock.On("UpdateLearnInterval", args.ctx, args.coll).Once().
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.args)
			err := wordService.UpdateLearnInterval(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
		})
	}
}

func Test_DeleteWord(t *testing.T) {
	ctx := context.Background()
	wordRepoMock, trMock := setupWordService(t)
	wordService := New(wordRepoMock, trMock)

	type args struct {
		ctx  context.Context
		coll entity.Collection
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(args args)
		wantErr   bool
	}{
		{
			name: "Add new word",
			args: args{
				ctx: ctx,
				coll: entity.Collection{
					Word:   "some_word",
					UserID: "12345",
					Name:   "some_coll",
				},
			},
			setupMock: func(args args) {
				wordRepoMock.On("DeleteWord", args.ctx, args.coll).Once().
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.args)
			err := wordService.DeleteWord(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
		})
	}
}
