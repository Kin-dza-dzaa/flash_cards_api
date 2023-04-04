package wordrepository

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/jackc/pgx/v4"
)

func (p *WordRepository) DeleteWord(ctx context.Context,
	collection entity.Collection,
) error {
	sql, args, err := p.Builder.Delete("*").
		From("user_collection").
		Where("user_id = ? AND word = ? AND collection_name = ?",
			collection.UserID, collection.Word, collection.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("WordRepository - DeleteWord - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("WordRepository - DeleteWord - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("WordRepository - DeleteWord - BeginFunc: %w", err)
	}

	return nil
}
