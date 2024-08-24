package slogz

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		slog.WarnContext(ctx, "context is nil")
		var pcs [1]uintptr
		// skip [runtime.Callers, this function]
		runtime.Callers(2, pcs[:])
		_ = slog.Default().Handler().Handle(ctx, slog.NewRecord(time.Now(), slog.LevelWarn, "context is nil", pcs[0]))
		return slog.Default()
	}

	l, ok := ctx.Value((*slog.Logger)(nil)).(*slog.Logger)
	if !ok {
		slog.WarnContext(ctx, "failed to get logger from context")
		return slog.Default()
	}

	return l
}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, (*slog.Logger)(nil), l)
}

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
