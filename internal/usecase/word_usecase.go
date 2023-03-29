// Package usecase implements application layer.
package usecase

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

type (
	database interface {
		IsWordInCollection(ctx context.Context, collection *entity.Collection) (bool, error)
		IsTransInDB(ctx context.Context, collection *entity.Collection) (bool, error)
		AddTranslation(ctx context.Context, wordTrans *entity.WordTrans) error
		AddWordToCollection(ctx context.Context, collection *entity.Collection) error
		UpdateLearnInterval(ctx context.Context, collection *entity.Collection) error
		DeleteWordFromCollection(ctx context.Context, collection *entity.Collection) error
		GetUserWords(ctx context.Context, collection *entity.Collection) (*entity.UserWords, error)
	}

	tranlsator interface {
		Translate(word string) (*entity.WordTrans, error)
	}
)

type WordService struct {
	dbAdapter          database
	googletransAdapter tranlsator
}

func (s *WordService) DeleteWordFromCollection(ctx context.Context,
	collection *entity.Collection) error {
	err := s.dbAdapter.DeleteWordFromCollection(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - DeleteWordFromCollection - "+
			"s.dbAdapter.DeleteWordFromCollection: %w", err)
	}
	return nil
}

func (s *WordService) UpdateLearnInterval(ctx context.Context,
	collection *entity.Collection) error {
	err := s.dbAdapter.UpdateLearnInterval(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - UpdateLearnInterval - "+
			"s.dbAdapter.UpdateLearnInterval: %w", err)
	}
	return nil
}

func (s *WordService) GetUserWords(ctx context.Context,
	collection *entity.Collection) (*entity.UserWords, error) {
	userWords, err := s.dbAdapter.GetUserWords(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("WordService - GetUserWords - s.dbAdapter.GetUserWords: %w", err)
	}
	return userWords, nil
}

func (s *WordService) AddWord(ctx context.Context, collection *entity.Collection) error {
	inCol, err := s.dbAdapter.IsWordInCollection(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - AddWord - s.dbAdapter.IsWordInCollection: %w", err)
	}
	if inCol {
		return nil
	}

	transInDB, err := s.dbAdapter.IsTransInDB(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - AddWord - s.dbAdapter.IsTransInDB: %w", err)
	}
	if !transInDB {
		if err := s.addTrans(ctx, collection.Word); err != nil {
			return fmt.Errorf("WordService - AddWord - s.addTrans: %w", err)
		}
	}

	err = s.dbAdapter.AddWordToCollection(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - AddWord - s.dbAdapter.AddTranslationToCollection: %w", err)
	}
	return nil
}

func (s *WordService) addTrans(ctx context.Context, word string) error {
	wordTrans, err := s.googletransAdapter.Translate(word)
	if err != nil {
		return fmt.Errorf("WordService - addTrans - s.googletransAdapter.Translate: %w", err)
	}

	return s.dbAdapter.AddTranslation(ctx, wordTrans)
}

func New(dbAdapter database, googletransAdapter tranlsator) *WordService {
	return &WordService{
		dbAdapter:          dbAdapter,
		googletransAdapter: googletransAdapter,
	}
}
