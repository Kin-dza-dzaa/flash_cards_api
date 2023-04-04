// Package wordservice implements application/use layer.
package wordservice

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

type (
	WordRepository interface {
		IsWordInCollection(ctx context.Context, collection entity.Collection) (bool, error)
		IsTransInDB(ctx context.Context, collection entity.Collection) (bool, error)
		AddTranslation(ctx context.Context, wordTrans entity.WordTrans) error
		AddWord(ctx context.Context, collection entity.Collection) error
		UpdateLearnInterval(ctx context.Context, collection entity.Collection) error
		DeleteWord(ctx context.Context, collection entity.Collection) error
		UserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error)
	}

	Translator interface {
		Translate(word string) (entity.WordTrans, error)
	}
)

type WordService struct {
	wordRepo         WordRepository
	googleTranslator Translator
}

func (s *WordService) DeleteWord(ctx context.Context,
	collection entity.Collection,
) error {
	err := s.wordRepo.DeleteWord(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - DeleteWord - s.wordRepo.DeleteWord: %w", err)
	}
	return nil
}

func (s *WordService) UpdateLearnInterval(ctx context.Context,
	collection entity.Collection,
) error {
	err := s.wordRepo.UpdateLearnInterval(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - UpdateLearnInterval - s.wordRepo.UpdateLearnInterval: %w", err)
	}
	return nil
}

func (s *WordService) UserWords(ctx context.Context,
	collection entity.Collection,
) (*entity.UserWords, error) {
	userWords, err := s.wordRepo.UserWords(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("WordService - UserWords - s.wordRepo.UserWords: %w", err)
	}
	return userWords, nil
}

func (s *WordService) AddWord(ctx context.Context, collection entity.Collection) error {
	inCol, err := s.wordRepo.IsWordInCollection(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - AddWord - s.wordRepo.IsWordInCollection: %w", err)
	}
	if inCol {
		return nil
	}

	transInDB, err := s.wordRepo.IsTransInDB(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - AddWord - s.wordRepo.IsTransInDB: %w", err)
	}
	if !transInDB {
		if err := s.addTrans(ctx, collection.Word); err != nil {
			return fmt.Errorf("WordService - AddWord - s.addTrans: %w", err)
		}
	}

	err = s.wordRepo.AddWord(ctx, collection)
	if err != nil {
		return fmt.Errorf("WordService - AddWord - s.wordRepo.AddWord: %w", err)
	}
	return nil
}

func (s *WordService) addTrans(ctx context.Context, word string) error {
	wordTrans, err := s.googleTranslator.Translate(word)
	if err != nil {
		return fmt.Errorf("WordService - addTrans - s.googleTranslator.Translate: %w", err)
	}

	return s.wordRepo.AddTranslation(ctx, wordTrans)
}

func New(dbAdapter WordRepository, googletransAdapter Translator) *WordService {
	return &WordService{
		wordRepo:         dbAdapter,
		googleTranslator: googletransAdapter,
	}
}
