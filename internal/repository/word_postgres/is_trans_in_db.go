package wordpostgres

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func (p *Postgres) IsTransInDB(ctx context.Context, collection entity.Collection) (bool, error) {
	subQuery := p.Builder.
		Select("*").
		From("word_translation").
		Where("word = ?", collection.Word)
	sql, args, err := sq.Expr("SELECT EXISTS(?)", subQuery).ToSql()
	if err != nil {
		return false, fmt.Errorf("Postgres - WordHasTranslation - ToSql: %w", err)
	}

	var transInDB bool
	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, sql, args...).Scan(&transInDB); err != nil {
			return fmt.Errorf("Postgres - WordHasTranslation - Scan: %w", err)
		}
		return nil
	})
	if err != nil {
		return transInDB, fmt.Errorf("Postgres - WordHasTranslation - BeginFunc: %w", err)
	}

	return transInDB, nil
}
