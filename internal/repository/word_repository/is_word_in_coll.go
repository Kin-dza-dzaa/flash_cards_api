package wordrepository

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func (p *WordRepository) IsWordInCollection(ctx context.Context,
	collection entity.Collection,
) (bool, error) {
	subQuery := p.Builder.
		Select("*").
		From("user_collection").
		Where("user_id = ? AND word = ? AND collection_name = ?",
			collection.UserID, collection.Word, collection.Name)
	sql, args, err := sq.Expr("SELECT EXISTS(?)", subQuery).ToSql()
	if err != nil {
		return false, fmt.Errorf("WordRepository - IsWordInCollection - ToSql: %w", err)
	}

	var inColl bool
	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, sql, args...).Scan(&inColl); err != nil {
			return fmt.Errorf("WordRepository - IsWordInCollection - Scan: %w", err)
		}
		return nil
	})
	if err != nil {
		return inColl, fmt.Errorf("WordRepository - IsWordInCollection - BeginFunc: %w", err)
	}

	return inColl, nil
}
