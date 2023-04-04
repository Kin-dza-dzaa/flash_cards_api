package wordrepository

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/google/go-cmp/cmp"
)

func setupIsTranslationInDB(ctx context.Context, t *testing.T, coll entity.Collection) *WordRepository {
	t.Helper()
	wordRepo := setupWordRepoContainer(t)

	if err := wordRepo.AddTranslation(ctx, entity.WordTrans{Word: coll.Word}); err != nil {
		t.Fatalf("setupIsWordInCollection - wordRepo.AddTranslation: %v", err)
	}

	return wordRepo
}

func Test_IsTranslationInDB(t *testing.T) {
	ctx := context.Background()
	existingColl := entity.Collection{
		Name:   "test_coll",
		Word:   "test_word",
		UserID: "12345",
	}
	wordRepo := setupIsTranslationInDB(ctx, t, existingColl)

	type args struct {
		ctx  context.Context
		coll entity.Collection
	}
	tests := []struct {
		name    string
		want    bool
		wantErr bool
		args    args
	}{
		{
			name: "Not existing trans",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "not_exist_word",
					UserID: "12345",
				},
				ctx: ctx,
			},
		},
		{
			name: "Existing trans",
			args: args{
				coll: existingColl,
				ctx:  ctx,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := wordRepo.IsTransInDB(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Fatalf("want %v but got: %v", tt.want, got)
			}
		})
	}
}
