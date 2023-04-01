package wordrepository

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func (p *WordRepository) IsTransInDB(ctx context.Context, collection entity.Collection) (bool, error) {
	subQuery := p.Builder.
		Select("*").
		From("word_translation").
		Where("word = ?", collection.Word)
	sql, args, err := sq.Expr("SELECT EXISTS(?)", subQuery).ToSql()
	if err != nil {
		return false, fmt.Errorf("WordRepository - IsTransInDB - ToSql: %w", err)
	}

	var transInDB bool
	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, sql, args...).Scan(&transInDB); err != nil {
			return fmt.Errorf("WordRepository - IsTransInDB - Scan: %w", err)
		}
		return nil
	})
	if err != nil {
		return transInDB, fmt.Errorf("WordRepository - IsTransInDB - BeginFunc: %w", err)
	}

	return transInDB, nil
}
