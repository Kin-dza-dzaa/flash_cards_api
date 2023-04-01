package wordrepository

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/jackc/pgx/v4"
)

func (p *WordRepository) AddTranslation(ctx context.Context, wordTrans entity.WordTrans) error {
	sql, args, err := p.Builder.
		Insert("word_translation").Columns("word, trans_data").
		Values(wordTrans.Word, wordTrans).
		ToSql()
	if err != nil {
		return fmt.Errorf("WordRepository - AddTranslation - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("WordRepository - AddTranslation - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("WordRepository - AddTranslation - BeginFunc: %w", err)
	}

	return nil
}