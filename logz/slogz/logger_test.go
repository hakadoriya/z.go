package slogz

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestRenewLogger(t *testing.T) {
	t.Parallel()
	t.Run("success,normal", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		logBuffer := new(bytes.Buffer)
		l := NewLogger(logBuffer, slog.LevelDebug, WithLoggerOptionAddCallerSkip(0), WithLoggerOptionHandlerOption(WithHandlerOptionAddSource(false), WithHandlerOptionAddTimestamp(false), WithHandlerOptionAddAttrs(slog.Bool("test", true))))
		l.Info("test", Error(io.EOF))
		t.Logf("logBuffer: %s", logBuffer.String())

		renewed := RenewLogger(ctx, l, WithLoggerOptionHandlerOption(WithHandlerOptionAddSource(false), WithHandlerOptionAddTimestamp(false), WithHandlerOptionAddAttrs(slog.Bool("test", true))))
		renewed.Info("test", Error(io.EOF))
		t.Logf("logBuffer: %s", logBuffer.String())
	})

	t.Run("failure,ErrHandlerIsNotSlogJSONHandler", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		logBuffer := new(bytes.Buffer)
		l := slog.New(slog.NewTextHandler(logBuffer, nil))
		_ = RenewLogger(ctx, l, WithLoggerOptionAddCallerSkip(0), WithLoggerOptionHandlerOption(WithHandlerOptionAddSource(false), WithHandlerOptionAddTimestamp(false), WithHandlerOptionAddAttrs(slog.Bool("test", true))))
		if !strings.Contains(logBuffer.String(), "handler=*slog.TextHandler: ") {
			t.Errorf("❌: logBuffer: %s", logBuffer.String())
		}
	})
}
