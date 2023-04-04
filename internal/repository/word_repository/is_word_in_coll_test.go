package wordrepository

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/google/go-cmp/cmp"
)

func setupIsWordInCollection(ctx context.Context, t *testing.T, coll entity.Collection) *WordRepository {
	t.Helper()
	wordRepo := setupWordRepoContainer(t)

	if err := wordRepo.AddTranslation(ctx, entity.WordTrans{Word: coll.Word}); err != nil {
		t.Fatalf("setupIsWordInCollection - wordRepo.AddTranslation: %v", err)
	}
	if err := wordRepo.AddWord(ctx, coll); err != nil {
		t.Fatalf("setupIsWordInCollection - wordRepo.AddWord: %v", err)
	}

	return wordRepo
}

func Test_IsWordInCollection(t *testing.T) {
	ctx := context.Background()
	existingWord := entity.Collection{
		Word:   "some_word",
		Name:   "test_coll",
		UserID: "12345",
	}
	wordRepo := setupIsWordInCollection(ctx, t, existingWord)

	type args struct {
		coll entity.Collection
		ctx  context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Not existing word",
			args: args{
				ctx: ctx,
			},
		},
		{
			name: "Existing word",
			args: args{
				coll: existingWord,
				ctx:  ctx,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := wordRepo.IsWordInCollection(tt.args.ctx, tt.args.coll)
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
