package wordrepository

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/google/go-cmp/cmp"
)

func Test_UserWords(t *testing.T) {
	ctx := context.Background()
	wordRepo := setupWordRepoContainer(t)

	type args struct {
		coll entity.Collection
		ctx  context.Context
	}
	tests := []struct {
		name          string
		wantUserWords *entity.UserWords
		wantErr       bool
		args          args
	}{
		{
			name: "Empty collection",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "test_word",
					UserID: "12345",
				},
				ctx: ctx,
			},
			wantUserWords: &entity.UserWords{
				Words: make(map[entity.CollectionName][]entity.WordData, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserWords, err := wordRepo.UserWords(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
			if diff := cmp.Diff(gotUserWords, tt.wantUserWords); diff != "" {
				t.Fatalf("user words must be equal diff: %v", diff)
			}
		})
	}
}
