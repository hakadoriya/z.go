package slogz

import (
	"fmt"
	"log/slog"

	"github.com/hakadoriya/z.go/pathz/filepathz"
)

func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case "source":
		switch v := a.Value.Any().(type) {
		case *slog.Source:
			return slog.String("file", filepathz.ExtractShortPath(fmt.Sprintf("%s:%d", v.File, v.Line)))
		default:
			return a
		}
	default:
		// no-op
	}

	return a
}
