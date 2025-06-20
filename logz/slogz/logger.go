package slogz

import (
	"context"
	"io"
	"log/slog"
)

func NewLogger(w io.Writer, level slog.Leveler, opts ...LoggerOption) *slog.Logger {
	c := new(slogJSONLoggerConfig)
	for _, o := range opts {
		o.apply(c)
	}
	handler := NewHandler(w, level, c.handlerOptions...)
	return slog.New(handler)
}

func RenewLogger(ctx context.Context, logger *slog.Logger, opts ...LoggerOption) *slog.Logger {
	c := new(slogJSONLoggerConfig)
	for _, o := range opts {
		o.apply(c)
	}
	handler := RenewHandler(ctx, logger.Handler(), c.handlerOptions...)
	return slog.New(handler)
}
