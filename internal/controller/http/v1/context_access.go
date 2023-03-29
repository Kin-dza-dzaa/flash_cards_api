package v1

import (
	"context"
	"fmt"
)

type key string

func inCtx(ctx context.Context, ctxKey string, value interface{}) context.Context {
	return context.WithValue(ctx, key(ctxKey), value)
}

func fromCtx[T interface{}](ctx context.Context, ctxKey string) (T, error) {
	val, ok := ctx.Value(key(ctxKey)).(T)
	if !ok {
		return val, fmt.Errorf("wordHandler - get - assertion %T:"+
			"couldn't get value from context", val)
	}
	return val, nil
}
