package slogz

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestFromContext(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		logBuffer := new(bytes.Buffer)
		l := FromContext(WithContext(context.Background(), slog.New(NewHandler(logBuffer, slog.LevelDebug).WithAttrs([]slog.Attr{slog.Bool("test", true)}))))
		l.Info("test")
		t.Logf("logBuffer: %s", logBuffer.String())
		requirez.StringHasPrefix(t, logBuffer.String(), `{"time":"`)
		requirez.StringContains(t, logBuffer.String(), `","severity":"INFO","caller":"`)
		requirez.StringContains(t, logBuffer.String(), `_test.go:`)
		requirez.StringHasSuffix(t, logBuffer.String(), `","msg":"test","test":true}`+"\n")
	})
}

//nolint:paralleltest
func TestFromContext_failure(t *testing.T) {
	t.Run("error,nil_context", func(t *testing.T) {
		logBuffer := new(bytes.Buffer)
		slog.SetDefault(slog.New(NewHandler(logBuffer, slog.LevelDebug).WithAttrs([]slog.Attr{slog.Bool("test", true)})))
		l := FromContext(nil)
		t.Logf("logBuffer: %s", logBuffer.String())
		requirez.Equal(t, slog.Default(), l)
		requirez.StringHasPrefix(t, logBuffer.String(), `{"time":"`)
		requirez.StringContains(t, logBuffer.String(), `","severity":"WARN","caller":"`)
		requirez.StringContains(t, logBuffer.String(), `_test.go:`)
		requirez.StringHasSuffix(t, logBuffer.String(), `","msg":"context is nil","test":true}`+"\n")
	})

	t.Run("error,not_found_in_context", func(t *testing.T) {
		logBuffer := new(bytes.Buffer)
		slog.SetDefault(slog.New(NewHandler(logBuffer, slog.LevelDebug).WithAttrs([]slog.Attr{slog.Bool("test", true)})))
		l := FromContext(context.Background())
		t.Logf("logBuffer: %s", logBuffer.String())
		requirez.Equal(t, slog.Default(), l)
		requirez.StringHasPrefix(t, logBuffer.String(), `{"time":"`)
		requirez.StringContains(t, logBuffer.String(), `","severity":"WARN","caller":"`)
		requirez.StringContains(t, logBuffer.String(), `_test.go:`)
		requirez.StringHasSuffix(t, logBuffer.String(), `","msg":"*slog.Logger not found in context","test":true}`+"\n")
	})
}
