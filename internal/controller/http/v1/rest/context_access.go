package rest

import (
	"context"
)

type key string

func inCtx(ctx context.Context, ctxKey string, value interface{}) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, key(ctxKey), value)
}

func fromCtx(ctx context.Context, ctxKey string) string {
	if ctx == nil {
		return ""
	}
	val, ok := ctx.Value(key(ctxKey)).(string)
	if !ok {
		return ""
	}
	return val
}
