package slogz

import (
	"log/slog"
)

// Error returns a slog.Attr with DefaultErrorKey and given error.
func Error(err error) slog.Attr {
	return slog.Any(DefaultErrorKey, err)
}
