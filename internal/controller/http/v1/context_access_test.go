package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fromCtx(t *testing.T) {
	Assert := assert.New(t)

	type args struct {
		Ctx    context.Context
		CtxKey string
	}
	tests := []struct {
		Name    string
		Args    args
		Want    string
		WantErr bool
	}{
		{
			Name:    "Not existing val",
			WantErr: true,
			Args: args{
				Ctx:    context.Background(),
				CtxKey: "key",
			},
		},
		{
			Name:    "Existing val",
			WantErr: false,
			Args: args{
				Ctx:    inCtx(context.Background(), "key", "val"),
				CtxKey: "key",
			},
			Want: "val",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := fromCtx[string](tc.Args.Ctx, tc.Args.CtxKey)
			if tc.WantErr {
				Assert.Error(err, "Err must be not nil")
			} else {
				Assert.Nil(err, "Err must be nil")
				Assert.Equal(tc.Want, got, "Val must be as expected")
			}
		})
	}
}
