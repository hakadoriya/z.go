package cliz

import (
	"context"
	"testing"

	"github.com/hakadoriya/z.go/cliz/clicorez"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_postCheckOptions(t *testing.T) {
	t.Parallel()

	t.Run("error,ErrDuplicateOption,Name", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		for _, subcmd := range c.SubCommands {
			subcmd.Options = append(subcmd.Options, &StringOption{Name: "foo"})
		}
		ctx := context.Background()
		_, err := c.parse(ctx, []string{"main-cli", "sub-cmd"})
		requirez.ErrorIs(t, err, clicorez.ErrDuplicateOption)
	})

	t.Run("error,ErrDuplicateOption,Aliases", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		for _, subcmd := range c.SubCommands {
			subcmd.Options = append(subcmd.Options, &StringOption{Name: "b2", Aliases: []string{"b"}})
		}
		ctx := context.Background()
		_, err := c.parse(ctx, []string{"main-cli", "sub-cmd"})
		requirez.ErrorIs(t, err, clicorez.ErrDuplicateOption)
	})

	t.Run("error,ErrOptionRequired", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		for _, subcmd := range c.SubCommands {
			subcmd.Options = append(subcmd.Options, &StringOption{Name: "required", Required: true})
		}
		ctx := context.Background()
		_, err := c.parse(ctx, []string{"main-cli", "sub-cmd"})
		requirez.ErrorIs(t, err, clicorez.ErrOptionRequired)
	})
}
