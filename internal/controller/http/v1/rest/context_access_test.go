package rest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_fromCtx(t *testing.T) {
	type args struct {
		Ctx    context.Context
		CtxKey string
	}
	tests := []struct {
		Name string
		Args args
		Want string
	}{
		{
			Name: "Not_existing_val",
			Args: args{
				Ctx:    context.Background(),
				CtxKey: "key",
			},
		},
		{
			Name: "Existing_val",
			Args: args{
				Ctx:    context.WithValue(context.Background(), key("key"), "val"),
				CtxKey: "key",
			},
			Want: "val",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := fromCtx(tt.Args.Ctx, tt.Args.CtxKey)
			if diff := cmp.Diff(tt.Want, got); diff != "" {
				t.Fatalf("wanted %v but got %v diff: %v", tt.Want, got, diff)
			}
		})
	}
}
