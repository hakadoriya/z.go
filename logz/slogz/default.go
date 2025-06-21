package slogz

import (
	"io"
	"log/slog"
	"os"
)

//nolint:gochecknoglobals
var (
	DefaultWriter     io.Writer    = os.Stdout
	DefaultLevel      slog.Leveler = slog.LevelDebug
	DefaultLevelKey                = "severity"
	DefaultSourceKey               = "caller"
	DefaultMessageKey              = "message"
	DefaultErrorKey                = "error"
)
