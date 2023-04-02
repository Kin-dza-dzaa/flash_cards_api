package wordrepository

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

func setupAddWord(ctx context.Context, coll entity.Collection, t *testing.T) *WordRepository {
	wordRepo := setupWordRepoContainer(ctx, t)

	if err := wordRepo.AddTranslation(ctx, entity.WordTrans{Word: coll.Word}); err != nil {
		t.Fatalf("setupAddWord - wordRepo.AddTranslation: %v", err)
	}

	return wordRepo
}

func Test_AddWord(t *testing.T) {
	ctx := context.Background()
	existingWord := entity.Collection{
		Name:   "test_coll",
		Word:   "test_word",
		UserID: "12345",
	}
	wordRepo := setupAddWord(ctx, existingWord, t)

	type args struct {
		coll entity.Collection
		ctx  context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add existing word",
			args: args{
				coll: existingWord,
				ctx:  ctx,
			},
		},
		{
			name: "Add not existing word",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "not_exist",
					UserID: "12345",
				},
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wordRepo.AddWord(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
		})
	}
}
