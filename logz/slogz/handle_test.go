package slogz

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func Test_slogHandler_WithGroup(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		logBuffer := new(bytes.Buffer)
		l := slog.New(NewHandler(logBuffer, slog.LevelDebug).WithAttrs([]slog.Attr{slog.Bool("test", true)}).WithGroup("testGroup"))
		l.Info("test", Error(io.EOF))
		t.Logf("logBuffer: %s", logBuffer.String())
		requirez.StringHasPrefix(t, logBuffer.String(), `{"time":"`)
		requirez.StringContains(t, logBuffer.String(), `","severity":"INFO","caller":"`)
		requirez.StringContains(t, logBuffer.String(), `_test.go:`)
		requirez.StringHasSuffix(t, logBuffer.String(), `","message":"test","test":true,"testGroup":{"error":"EOF","errorVerbose":"EOF"}}`+"\n")
	})
}
