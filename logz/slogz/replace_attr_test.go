package slogz

import (
	"log/slog"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func TestReplaceAttr(t *testing.T) {
	t.Parallel()

	t.Run("success,severity", func(t *testing.T) {
		t.Parallel()

		actual := ReplaceAttr(nil, slog.Any("level", slog.LevelDebug))
		requirez.Equal(t, `severity=DEBUG`, actual.String())
	})
	t.Run("success,not_severity", func(t *testing.T) {
		t.Parallel()

		actual := ReplaceAttr(nil, slog.Any("level", "100"))
		requirez.Equal(t, `level=100`, actual.String())
	})
	t.Run("success,not_caller", func(t *testing.T) {
		t.Parallel()

		actual := ReplaceAttr(nil, slog.String("source", "SOURCE"))
		requirez.Equal(t, `source=SOURCE`, actual.String())
	})
}
