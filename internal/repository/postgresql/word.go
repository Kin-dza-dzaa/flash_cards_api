// Package postgresql implements adapter layer for postgres database.
package postgresql

import (
	"context"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/service"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"go.opentelemetry.io/otel"
)

const otelName = "github.com/Kin-dza-dzaa/flash_cards_api/internal/repository/postgresql"

var _ = service.WordRepo((*Word)(nil))

type Word struct {
	*postgres.ConnPool
}

func (p *Word) Close() {
	p.Pool.Close()
}

func (p *Word) UserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "WordPostgresql - UserWords")
	defer span.End()

	sql, args, err := p.Builder.Select("collection_name, time_diff, last_repeat, trans_data").
		From("user_collection").
		Join("word_translation USING(word)").
		Where("user_id = ?", collection.UserID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Word - UserWords - ToSql: %w", err)
	}

	userWords := new(entity.UserWords)
	userWords.Words = make(map[entity.CollectionName][]entity.WordData)
	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Word - UserWords - Query: %w", err)
		}
		defer rows.Close()

		collectionName, wordData := entity.CollectionName(""), entity.WordData{}
		for rows.Next() {
			if err := rows.Scan(&collectionName, &wordData.TimeDiff, &wordData.LastRepeat, &wordData.WordTrans); err != nil {
				return fmt.Errorf("Word - UserWords - Scan: %w", err)
			}

			userWords.Words[collectionName] = append(userWords.Words[collectionName], wordData)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Word - UserWords - BeginFunc: %w", err)
	}

	return userWords, nil
}

func (p *Word) UpdateLearnInterval(ctx context.Context, collection entity.Collection) error {
	_, span := otel.Tracer(otelName).Start(ctx, "WordPostgresql - UpdateLearnInterval")
	defer span.End()

	sql, args, err := p.Builder.Update("user_collection").
		Set("time_diff", collection.TimeDiff).
		Set("last_repeat", collection.LastRepeat).
		Where("user_id = ? AND word = ? AND collection_name = ?",
			collection.UserID, collection.Word, collection.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("Word - UpdateLearnInterval - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Word - UpdateLearnInterval - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Word - UpdateLearnInterval - BeginFunc: %w", err)
	}

	return nil
}

func (p *Word) IsWordInCollection(ctx context.Context, collection entity.Collection) (bool, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "WordPostgresql - IsWordInCollection")
	defer span.End()

	subQuery := p.Builder.
		Select("*").
		From("user_collection").
		Where("user_id = ? AND word = ? AND collection_name = ?",
			collection.UserID, collection.Word, collection.Name)
	sql, args, err := sq.Expr("SELECT EXISTS(?)", subQuery).ToSql()
	if err != nil {
		return false, fmt.Errorf("Word - IsWordInCollection - ToSql: %w", err)
	}

	var inColl bool
	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, sql, args...).Scan(&inColl); err != nil {
			return fmt.Errorf("Word - IsWordInCollection - Scan: %w", err)
		}
		return nil
	})
	if err != nil {
		return inColl, fmt.Errorf("Word - IsWordInCollection - BeginFunc: %w", err)
	}

	return inColl, nil
}

func (p *Word) IsTransInDB(ctx context.Context, collection entity.Collection) (bool, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "WordPostgresql - IsTransInDB")
	defer span.End()

	subQuery := p.Builder.
		Select("*").
		From("word_translation").
		Where("word = ?", collection.Word)
	sql, args, err := sq.Expr("SELECT EXISTS(?)", subQuery).ToSql()
	if err != nil {
		return false, fmt.Errorf("Word - IsTransInDB - ToSql: %w", err)
	}

	var transInDB bool
	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, sql, args...).Scan(&transInDB); err != nil {
			return fmt.Errorf("Word - IsTransInDB - Scan: %w", err)
		}
		return nil
	})
	if err != nil {
		return transInDB, fmt.Errorf("Word - IsTransInDB - BeginFunc: %w", err)
	}

	return transInDB, nil
}

func (p *Word) DeleteWord(ctx context.Context, collection entity.Collection) error {
	_, span := otel.Tracer(otelName).Start(ctx, "WordPostgresql - DeleteWord")
	defer span.End()

	sql, args, err := p.Builder.Delete("*").
		From("user_collection").
		Where("user_id = ? AND word = ? AND collection_name = ?",
			collection.UserID, collection.Word, collection.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("Word - DeleteWord - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Word - DeleteWord - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Word - DeleteWord - BeginFunc: %w", err)
	}

	return nil
}

func (p *Word) AddWord(ctx context.Context, collection entity.Collection) error {
	_, span := otel.Tracer(otelName).Start(ctx, "WordPostgresql - AddWord")
	defer span.End()

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
		return fmt.Errorf("Word - AddTranslation - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Word - AddTranslation - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Word - AddTranslation - BeginFunc: %w", err)
	}

	return nil
}

func (p *Word) AddTranslation(ctx context.Context, wordTrans entity.WordTrans) error {
	_, span := otel.Tracer(otelName).Start(ctx, "WordPostgresql - AddTranslation")
	defer span.End()

	sql, args, err := p.Builder.
		Insert("word_translation").Columns("word, trans_data").
		Values(wordTrans.Word, wordTrans).
		ToSql()
	if err != nil {
		return fmt.Errorf("Word - AddTranslation - ToSql: %w", err)
	}

	err = p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("Word - AddTranslation - Exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Word - AddTranslation - BeginFunc: %w", err)
	}

	return nil
}

func NewWordPostgre(pool *postgres.ConnPool) *Word {
	return &Word{
		pool,
	}
}
