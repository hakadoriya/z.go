package cliz

import (
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func TestCommand_loadDefaults(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		for _, subcmd := range c.SubCommands {
			subcmd.Options = append(c.Options, &unknownOptionType{Name: "unknown", Environment: "UNKNOWN", Description: "my unknown option"})
		}
		err := c.loadDefaults()
		requirez.ErrorIs(t, err, ErrInvalidOptionType)
	})
}
