package slogz

import (
	"fmt"
	"log/slog"

	"github.com/hakadoriya/z.go/pathz/filepathz"
)

func ReplaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case "level":
		_, ok := a.Value.Any().(slog.Level)
		if !ok {
			return a
		}

		a.Key = "severity"
		return a
	case "source":
		switch v := a.Value.Any().(type) {
		case *slog.Source:
			return slog.String("caller", filepathz.ExtractShortPath(fmt.Sprintf("%s:%d", v.File, v.Line)))
		default:
			return a
		}
	default:
		// no-op
	}

	return a
}
