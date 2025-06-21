package slogz

import "context"

type ctxKeyAddCallerSkip struct{}

func ContextWithAddCallerSkip(ctx context.Context, skip int) context.Context {
	currentSkip := contextAddCallerSkip(ctx)

	return context.WithValue(ctx, ctxKeyAddCallerSkip{}, currentSkip+skip)
}

func contextAddCallerSkip(ctx context.Context) int {
	if ctx == nil {
		return 0
	}

	skip, ok := ctx.Value(ctxKeyAddCallerSkip{}).(int)
	if !ok {
		return 0
	}

	return skip
}
