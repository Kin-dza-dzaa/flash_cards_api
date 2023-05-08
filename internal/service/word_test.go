package service

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/service/repomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
)

func setupWordService(t *testing.T) (*repomock.WordRepo, *repomock.TransRepo) {
	t.Helper()
	db := repomock.NewWordRepo(t)
	tr := repomock.NewTransRepo(t)

	return db, tr
}

func Test_AddWord(t *testing.T) {
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args)
		wantErr   bool
	}{
		{
			name: "Add new word",
			args: args{
				coll: entity.Collection{
					Name: "some_name",
					Word: "Some_words",
				},
			},
			setupMock: func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args) {
				dbMock.On("IsWordInCollection", mock.Anything, args.coll).Once().Return(false, nil)
				dbMock.On("IsTransInDB", mock.Anything, args.coll).Once().Return(false, nil)
				trMock.On("Translate", mock.Anything, args.coll.Word).Once().
					Return(entity.WordTrans{}, nil)
				dbMock.On("AddTranslation", mock.Anything, mock.Anything).Once().
					Return(nil)
				dbMock.On("AddWord", mock.Anything, args.coll).Once().Return(nil)
			},
		},
		{
			name: "Add existing word",
			args: args{
				coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			setupMock: func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args) {
				dbMock.On("IsWordInCollection", mock.Anything, args.coll).Once().Return(true, nil)
			},
		},
		{
			name: "Add word that in DB",
			args: args{
				coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			setupMock: func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args) {
				dbMock.On("IsWordInCollection", mock.Anything, args.coll).Once().Return(false, nil)
				dbMock.On("IsTransInDB", mock.Anything, args.coll).Once().Return(true, nil)
				dbMock.On("AddWord", mock.Anything, args.coll).Once().Return(nil)
			},
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		dbMock, trMock := setupWordService(t)
		wordService := NewWordService(dbMock, trMock)
		tt.setupMock(dbMock, trMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			err := wordService.AddWord(ctx, tt.args.coll)
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
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name          string
		args          args
		setupMock     func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args)
		wantErr       bool
		wantUserWords *entity.UserWords
	}{
		{
			name: "Add new word",
			args: args{
				coll: entity.Collection{
					Name:   "some_name",
					UserID: "12345",
					Word:   "Some_words",
				},
			},
			setupMock: func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args) {
				dbMock.On("UserWords", mock.Anything, args.coll).Once().
					Return(new(entity.UserWords), nil)
			},
			wantErr:       false,
			wantUserWords: new(entity.UserWords),
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		dbMock, trMock := setupWordService(t)
		wordService := NewWordService(dbMock, trMock)
		tt.setupMock(dbMock, trMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			gotUserWords, err := wordService.UserWords(ctx, tt.args.coll)
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
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args)
		wantErr   bool
	}{
		{
			name: "Add new word",
			args: args{
				coll: entity.Collection{
					Word:   "some_word",
					UserID: "12345",
					Name:   "some_coll",
				},
			},
			setupMock: func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args) {
				dbMock.On("UpdateLearnInterval", mock.Anything, args.coll).Once().
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		dbMock, trMock := setupWordService(t)
		wordService := NewWordService(dbMock, trMock)
		tt.setupMock(dbMock, trMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			err := wordService.UpdateLearnInterval(ctx, tt.args.coll)
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
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args)
		wantErr   bool
	}{
		{
			name: "Add new word",
			args: args{
				coll: entity.Collection{
					Word:   "some_word",
					UserID: "12345",
					Name:   "some_coll",
				},
			},
			setupMock: func(dbMock *repomock.WordRepo, trMock *repomock.TransRepo, args args) {
				dbMock.On("DeleteWord", mock.Anything, args.coll).Once().
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		dbMock, trMock := setupWordService(t)
		wordService := NewWordService(dbMock, trMock)
		tt.setupMock(dbMock, trMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			err := wordService.DeleteWord(ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
		})
	}
}
