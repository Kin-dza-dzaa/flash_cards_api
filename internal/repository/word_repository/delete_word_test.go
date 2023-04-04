package wordrepository

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

func Test_DeleteWord(t *testing.T) {
	ctx := context.Background()
	wordRepo := setupWordRepoContainer(t)

	type args struct {
		ctx  context.Context
		coll entity.Collection
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete word",
			args: args{
				coll: entity.Collection{
					Name:   "test_coll",
					Word:   "test_word",
					UserID: "12345",
				},
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Delete not existing word",
			args: args{
				ctx: ctx,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wordRepo.DeleteWord(tt.args.ctx, tt.args.coll)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
		})
	}
}
