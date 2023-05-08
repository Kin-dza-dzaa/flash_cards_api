package postgresql

import (
	"context"
	"fmt"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	"github.com/adrianbrad/psqldocker"
	"github.com/google/go-cmp/cmp"
)

func Test_UserWords(t *testing.T) {
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name          string
		wantUserWords *entity.UserWords
		wantErr       bool
		args          args
	}{
		{
			name: "Empty_collection",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "test_word",
					UserID: "12345",
				},
			},
			wantUserWords: &entity.UserWords{
				Words: make(map[entity.CollectionName][]entity.WordData, 0),
			},
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		wordRepo := setupWordRepoContainer(ctx, t, tt.name)

		t.Run(tt.name, func(t *testing.T) {
			gotUserWords, err := wordRepo.UserWords(ctx, tt.args.coll)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err but got: %v", err)
			}
			if diff := cmp.Diff(gotUserWords, tt.wantUserWords); diff != "" {
				t.Fatalf("user words must be equal diff: %v", diff)
			}
		})
	}
}

func Test_UpdateLearnInterval(t *testing.T) {
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update_interval",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "test_word",
					UserID: "12345",
				},
			},
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		wordRepo := setupWordRepoContainer(ctx, t, tt.name)

		t.Run(tt.name, func(t *testing.T) {
			err := wordRepo.UpdateLearnInterval(ctx, tt.args.coll)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err but got: %v", err)
			}
		})
	}
}

func Test_IsWordInCollection(t *testing.T) {
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Not_existing_word",
			args: args{
				coll: entity.Collection{
					Word:   "not_existing_word",
					Name:   "test_coll",
					UserID: "12345",
				},
			},
		},
		{
			name: "Existing_word",
			args: args{
				coll: entity.Collection{
					Word:   "some_word",
					Name:   "test_coll",
					UserID: "12345",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		wordRepo := setupWordRepoContainer(ctx, t, tt.name)
		if tt.want {
			setupAddTranslationToDB(ctx, t, tt.args.coll, wordRepo)
			setupAddWordToUser(ctx, t, tt.args.coll, wordRepo)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := wordRepo.IsWordInCollection(ctx, tt.args.coll)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err but got: %v", err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Fatalf("want %v but got: %v", tt.want, got)
			}
		})
	}
}

func Test_IsTranslationInDB(t *testing.T) {
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name    string
		want    bool
		wantErr bool
		args    args
	}{
		{
			name: "Not_existing_trans",
			args: args{
				coll: entity.Collection{
					Name:   "not existing trans",
					Word:   "not_exist_word",
					UserID: "12345",
				},
			},
		},
		{
			name: "Existing_trans",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "test_word",
					UserID: "12345",
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		wordRepo := setupWordRepoContainer(ctx, t, tt.name)
		if tt.want {
			setupAddTranslationToDB(ctx, t, tt.args.coll, wordRepo)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := wordRepo.IsTransInDB(ctx, tt.args.coll)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err but got: %v", err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Fatalf("want %v but got: %v", tt.want, got)
			}
		})
	}
}

func Test_DeleteWord(t *testing.T) {
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete_word",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "test_word",
					UserID: "12345",
				},
			},
			wantErr: false,
		},
		{
			name: "Delete_not_existing_word",
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		wordRepo := setupWordRepoContainer(ctx, t, tt.name)

		t.Run(tt.name, func(t *testing.T) {
			err := wordRepo.DeleteWord(ctx, tt.args.coll)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err but got: %v", err)
			}
		})
	}
}

func Test_AddWord(t *testing.T) {
	type args struct {
		coll entity.Collection
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add_existing_word",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "test_word",
					UserID: "12345",
				},
			},
		},
		{
			name: "Add_not_existing_word",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "not_exist",
					UserID: "12345",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		wordRepo := setupWordRepoContainer(ctx, t, tt.name)
		if !tt.wantErr {
			setupAddTranslationToDB(ctx, t, tt.args.coll, wordRepo)
		}

		t.Run(tt.name, func(t *testing.T) {
			err := wordRepo.AddWord(ctx, tt.args.coll)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err but got: %v", err)
			}
		})
	}
}

func Test_AddTrans(t *testing.T) {
	type args struct {
		wordTrans entity.WordTrans
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add_normal_word",
			args: args{
				wordTrans: entity.WordTrans{
					Word: "test_word",
				},
			},
			wantErr: false,
		},
		{
			name:    "Add_empty_word",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		wordRepo := setupWordRepoContainer(ctx, t, tt.name)

		t.Run(tt.name, func(t *testing.T) {
			err := wordRepo.AddTranslation(ctx, tt.args.wordTrans)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err but got: %v", err)
			}
		})
	}
}

func setupAddTranslationToDB(ctx context.Context, t *testing.T, coll entity.Collection, wordRepo *Word) {
	t.Helper()

	// Add translation to DB.
	sql, args, err := wordRepo.Builder.
		Insert("word_translation").Columns("word, trans_data").
		Values(coll.Word, entity.WordTrans{}).
		ToSql()
	if err != nil {
		t.Fatalf("wordRepo.Builder.ToSql: %v", err)
	}
	if _, err := wordRepo.Pool.Exec(ctx, sql, args...); err != nil {
		t.Fatalf("add translation failed: %v", err)
	}
}

func setupAddWordToUser(ctx context.Context, t *testing.T, coll entity.Collection, wordRepo *Word) {
	t.Helper()

	sql, args, err := wordRepo.Builder.
		Insert("user_collection").
		Columns("user_id, word, collection_name, time_diff, last_repeat").
		Values(
			coll.UserID,
			coll.Word,
			coll.Name,
			coll.TimeDiff,
			coll.LastRepeat,
		).
		ToSql()
	if err != nil {
		t.Fatalf("wordRepo.Builder.ToSql: %v", err)
	}
	if _, err := wordRepo.Pool.Exec(ctx, sql, args...); err != nil {
		t.Fatalf("add word to user collection failed: %v", err)
	}
}

// Creates new throw away postgres:alpine container.
func setupWordRepoContainer(ctx context.Context, t *testing.T, containerName string) *Word {
	t.Helper()
	const (
		db   = "test_db"
		user = "user"
		pass = "password"
	)

	sql := `
		CREATE TABLE IF NOT EXISTS word_translation(
			word                TEXT                                                                NOT NULL CHECK(word != ''),
			trans_data          JSONB                                                               NOT NULL,
			PRIMARY KEY (word)
		);
		
		CREATE TABLE IF NOT EXISTS user_collection(
			user_id                                     TEXT                                        NOT NULL,
			word                                        TEXT                                        NOT NULL CHECK(word != ''),
			collection_name                             TEXT                                        NOT NULL,
			time_diff                                   INTERVAL                                    NOT NULL,
			last_repeat                                 TIMESTAMP                                   NOT NULL,
			FOREIGN KEY (word) REFERENCES word_translation(word),
			UNIQUE(user_id, word, collection_name)
		);	
	`

	t.Log("starting up a psql container")
	c, err := psqldocker.NewContainer(
		user,
		pass,
		db,
		psqldocker.WithContainerName(containerName),
		psqldocker.WithSQL(sql),
	)
	if err != nil {
		t.Fatalf("setupWordRepo - psqldocker.NewContainer: %v", err)
	}
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Fatalf("setupWordRepo - Cleanup - c.Close: %v", err)
		}
	})

	connPool, err := postgres.New(ctx, fmt.Sprintf("postgresql://%s:%s@0.0.0.0:%s/%s", user, pass, c.Port(), db), 10)
	if err != nil {
		t.Fatalf("setupWordRepo - postgres.New: %v", err)
	}

	return NewWordPostgre(connPool)
}
