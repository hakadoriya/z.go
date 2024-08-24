package cliz

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/hakadoriya/z.go/cliz/clierrz"
	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("success,ExecFunc", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		ctx := context.Background()
		err := c.Exec(ctx, []string{"main-cli", "sub-cmd"})
		requirez.NoError(t, err)
	})

	t.Run("success,AllHookExecFunc", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		ctx := context.Background()
		err := c.Exec(ctx, []string{"main-cli", "sub-cmd", "sub-sub-cmd", "--id", "1"})
		requirez.NoError(t, err)
	})

	t.Run("error,c.Parse,ErrCommandExecIsNotDefined", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		ctx := context.Background()
		err := c.Exec(ctx, []string{"main-cli", "--unknown"})
		requirez.ErrorIs(t, err, clierrz.ErrUnknownOption)
	})

	t.Run("error,PreHookExecFunc", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		ctx := context.Background()
		err := c.Exec(ctx, []string{"main-cli", "sub-cmd", "sub-sub-cmd2", "--id", "1"})
		requirez.ErrorIs(t, err, io.ErrUnexpectedEOF)
	})

	t.Run("error,ErrCommandExecIsNotDefined", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		ctx := context.Background()
		buf := new(bytes.Buffer)
		c.SetStdout(buf)
		c.SetStderr(buf)
		err := c.Exec(ctx, []string{"main-cli"})
		requirez.ErrorIs(t, err, clierrz.ErrHelp)
		assertz.StringContains(t, buf.String(), c.Description)
	})

	t.Run("error,ExecFunc", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		ctx := context.Background()
		err := c.Exec(ctx, []string{"main-cli", "sub-cmd", "sub-sub-cmd3", "--id", "1"})
		requirez.ErrorIs(t, err, io.ErrUnexpectedEOF)
	})

	t.Run("error,PostHookExecFunc", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		ctx := context.Background()
		err := c.Exec(ctx, []string{"main-cli", "sub-cmd", "sub-sub-cmd4", "--id", "1"})
		requirez.ErrorIs(t, err, io.ErrUnexpectedEOF)
	})
}
