package slogz

import (
	"log/slog"

	"github.com/hakadoriya/z.go/logz/slogz/slogcorez"
)

// Error returns a slog.Attr with DefaultErrorKey and given error.
func Error(err error) slog.Attr {
	return slog.Any(slogcorez.DefaultErrorKey, err)
}
