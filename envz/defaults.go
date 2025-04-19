package envz

import (
	"log/slog"
	"os"

	"github.com/hakadoriya/z.go/logz/slogz"
)

//nolint:gochecknoglobals
var (
	DefaultTagKey      = "env"
	DefaultRequiredKey = "required"
	DefaultDefaultKey  = "default"
	Logger             = slog.New(slogz.NewHandler(os.Stdout, slog.LevelInfo))
	// Logger = slog.New(slogz.NewHandler(os.Stdout, slog.LevelDebug))
)
