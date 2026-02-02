package app

import "context"

type contextKey string

const (
	ContextUserKey contextKey = "user"
)

func ContextSetUser(ctx context.Context, u any) context.Context {
	return context.WithValue(ctx, ContextUserKey, u)
}

func ContextGetUser(ctx context.Context) (any, bool) {
	v := ctx.Value(ContextUserKey)
	if v == nil {
		return nil, false
	}
	return v, true
}
