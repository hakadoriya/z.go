package slogz

import (
	"context"
	"log/slog"
)

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		//nolint:contextcheck // avoid panic
		slog.WarnContext(ContextWithAddCallerSkip(context.Background(), 1), "context is nil")
		return slog.Default()
	}

	l, ok := ctx.Value((*slog.Logger)(nil)).(*slog.Logger)
	if !ok {
		slog.WarnContext(ContextWithAddCallerSkip(ctx, 1), "*slog.Logger not found in context")
		return slog.Default()
	}

	return l
}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, (*slog.Logger)(nil), l)
}
