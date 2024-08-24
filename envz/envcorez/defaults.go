package envcorez

import (
	"log/slog"
	"os"

	"github.com/hakadoriya/z.go/logz/slogz"
)

//nolint:gochecknoglobals
var (
	TagKey            = "env"
	OptionKeyRequired = "required"
	OptionKeyDefault  = "default"
	Logger            = slog.New(slogz.NewHandler(os.Stdout, slog.LevelDebug))
)
