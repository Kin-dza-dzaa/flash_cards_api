package wordpostgres

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/jackc/pgx/v4"
)

func (p *Postgres) UpdateLearnInterval(ctx context.Context, collection entity.Collection) error {
	sql, args, err := p.Builder.Update("user_collection").
		Set("time_diff", collection.TimeDiff).
		Set("last_repeat", collection.LastRepeat).
		Where("user_id = ? AND word = ? AND collection_name = ?",
			collection.UserID, collection.Word, collection.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("Postgres - UpdateLearnInterval - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Postgres - UpdateLearnInterval - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Postgres - UpdateLearnInterval - BeginFunc: %w", err)
	}

	return nil
}
