package wordpostgres

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/jackc/pgx/v4"
)

func (p *Postgres) GetUserWords(ctx context.Context,
	collection *entity.Collection) (*entity.UserWords, error) {
	sql, args, err := p.Builder.Select("collection_name, time_diff, last_repeat, trans_data").
		From("user_collection").
		Join("word_translation USING(word)").
		Where("user_id = ?", collection.UserID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Postgres - GetUserWords - ToSql: %w", err)
	}

	var userWords = new(entity.UserWords)
	userWords.Words = make(map[entity.CollectionName][]entity.WordData)
	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Postgres - GetUserWords - Query: %w", err)
		}
		defer rows.Close()

		collectionName, wordData := entity.CollectionName(""), entity.WordData{}
		for rows.Next() {
			if err := rows.Scan(&collectionName, &wordData.TimeDiff, &wordData.LastRepeat, &wordData.WordTrans); err != nil {
				return fmt.Errorf("Postgres - GetUserWords - Scan: %w", err)
			}

			userWords.Words[collectionName] = append(userWords.Words[collectionName], wordData)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Postgres - GetUserWords - BeginFunc: %w", err)
	}

	return userWords, nil
}
