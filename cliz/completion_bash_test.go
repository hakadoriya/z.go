package cliz

import (
	"bytes"
	"context"
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_initAppendGenerateBashCompletionSubCommands(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		buf := new(bytes.Buffer)
		c.SetStdout(buf)
		ctx := context.Background()
		err := c.Exec(ctx, []string{"main-cli", "sub-cmd", DefaultGenerateBashCompletionSubCommandName})
		requirez.NoError(t, err)
	})
}
