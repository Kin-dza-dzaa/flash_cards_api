package wordpostgres

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/jackc/pgx/v4"
)

func (p *Postgres) AddWordToCollection(ctx context.Context,
	collection entity.Collection) error {
	sql, args, err := p.Builder.Insert("user_collection").
		Columns("user_id, word, collection_name, time_diff, last_repeat").
		Values(
			collection.UserID,
			collection.Word,
			collection.Name,
			collection.TimeDiff,
			collection.LastRepeat,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("Postgres - AddTranslationToCollection - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Postgres - AddTranslationToCollection - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Postgres - AddTranslationToCollection - BeginFunc: %w", err)
	}

	return nil
}
