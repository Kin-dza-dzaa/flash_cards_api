package wordrepository

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

func Test_AddTrans(t *testing.T) {
	ctx := context.Background()
	wordRepo := setupWordRepoContainer(ctx, t)

	type args struct {
		wordTrans entity.WordTrans
		ctx       context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add normal word",
			args: args{
				wordTrans: entity.WordTrans{
					Word: "test_word",
				},
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Add empty word",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wordRepo.AddTranslation(tt.args.ctx, tt.args.wordTrans)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
		})
	}
}
