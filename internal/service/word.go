// Package service implements use-case/application layer, one service per file.
package service

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"go.opentelemetry.io/otel"
)

const otelName = "github.com/Kin-dza-dzaa/flash_cards_api/internal/service"

type (
	WordRepo interface {
		IsWordInCollection(ctx context.Context, collection entity.Collection) (bool, error)
		IsTransInDB(ctx context.Context, collection entity.Collection) (bool, error)
		AddTranslation(ctx context.Context, wordTrans entity.WordTrans) error
		AddWord(ctx context.Context, collection entity.Collection) error
		UpdateLearnInterval(ctx context.Context, collection entity.Collection) error
		DeleteWord(ctx context.Context, collection entity.Collection) error
		UserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error)
	}

	TransRepo interface {
		Translate(ctx context.Context, word string) (entity.WordTrans, error)
	}
)

type Word struct {
	wordRepo  WordRepo
	transRepo TransRepo
}

func (s *Word) DeleteWord(ctx context.Context, collection entity.Collection) error {
	_, span := otel.Tracer(otelName).Start(ctx, "WordService - DeleteWord")
	defer span.End()

	err := s.wordRepo.DeleteWord(ctx, collection)
	if err != nil {
		return fmt.Errorf("Word - DeleteWord - s.wordRepo.DeleteWord: %w", err)
	}
	return nil
}

func (s *Word) UpdateLearnInterval(ctx context.Context, collection entity.Collection) error {
	_, span := otel.Tracer(otelName).Start(ctx, "WordService - UpdateLearnInterval")
	defer span.End()

	err := s.wordRepo.UpdateLearnInterval(ctx, collection)
	if err != nil {
		return fmt.Errorf("Word - UpdateLearnInterval - s.wordRepo.UpdateLearnInterval: %w", err)
	}
	return nil
}

func (s *Word) UserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "WordService - UserWords")
	defer span.End()

	userWords, err := s.wordRepo.UserWords(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("Word - UserWords - s.wordRepo.UserWords: %w", err)
	}
	return userWords, nil
}

func (s *Word) AddWord(ctx context.Context, collection entity.Collection) error {
	_, span := otel.Tracer(otelName).Start(ctx, "WordService - AddWord")
	defer span.End()

	inCol, err := s.wordRepo.IsWordInCollection(ctx, collection)
	if err != nil {
		return fmt.Errorf("Word - AddWord - s.wordRepo.IsWordInCollection: %w", err)
	}
	if inCol {
		return nil
	}

	transInDB, err := s.wordRepo.IsTransInDB(ctx, collection)
	if err != nil {
		return fmt.Errorf("Word - AddWord - s.wordRepo.IsTransInDB: %w", err)
	}
	if !transInDB {
		if err := s.addTrans(ctx, collection.Word); err != nil {
			return fmt.Errorf("Word - AddWord - s.addTrans: %w", err)
		}
	}

	err = s.wordRepo.AddWord(ctx, collection)
	if err != nil {
		return fmt.Errorf("Word - AddWord - s.wordRepo.AddWord: %w", err)
	}
	return nil
}

func (s *Word) addTrans(ctx context.Context, word string) error {
	wordTrans, err := s.transRepo.Translate(ctx, word)
	if err != nil {
		return fmt.Errorf("Word - addTrans - s.googleTranslator.Translate: %w", err)
	}

	return s.wordRepo.AddTranslation(ctx, wordTrans)
}

func NewWordService(wordRepo WordRepo, translatorRepo TransRepo) *Word {
	return &Word{
		wordRepo:  wordRepo,
		transRepo: translatorRepo,
	}
}
