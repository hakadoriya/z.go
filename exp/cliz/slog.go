package cliz

import (
	"io"
	"log/slog"

	"github.com/hakadoriya/z.go/logz/slogz"
)

var (
	Logger = slog.New(slogz.NewHandler(io.Discard, slog.LevelDebug))
)
