package slogz

import "context"

type ctxKeyAddCallerSkip struct{}

func AddCallerSkip(ctx context.Context, skip int) context.Context {
	return context.WithValue(ctx, ctxKeyAddCallerSkip{}, skip)
}

func addCallerSkip(ctx context.Context) int {
	if ctx == nil {
		return 0
	}

	skip, ok := ctx.Value(ctxKeyAddCallerSkip{}).(int)
	if !ok {
		return 0
	}

	return skip
}
