package cliz

import (
	"io"
	"log/slog"
	"os"

	"github.com/hakadoriya/z.go/logz/slogz"
)

//nolint:gochecknoglobals
var (
	DefaultCompletionSubCommandName             = "completion"
	DefaultGenerateBashCompletionSubCommandName = "__generate_bash_completion"
	DefaultGenerateZshCompletionSubCommandName  = "__generate_zsh_completion"

	DefaultTagKey         = "cli"
	DefaultAliasKey       = "alias"
	DefaultEnvKey         = "env"
	DefaultDefaultKey     = "default"
	DefaultRequiredKey    = "required"
	DefaultDescriptionKey = "description"
	DefaultHiddenKey      = "hidden"

	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr

	Logger = slog.New(slogz.NewHandler(io.Discard, slog.LevelDebug))
)
