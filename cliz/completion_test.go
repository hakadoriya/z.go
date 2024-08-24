package cliz

import (
	"bytes"
	"context"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

//nolint:paralleltest
func TestCommand_initAppendCompletionSubCommand(t *testing.T) {
	t.Run("success,", func(t *testing.T) {
		c := newTestCommand()
		buf := new(bytes.Buffer)
		c.SetStdout(buf)
		c.SetStderr(buf)
		err := c.Exec(context.Background(), []string{"main-cli", "completion", "bash"})
		requirez.NoError(t, err)
		assertz.StringContains(t, buf.String(), c.Name)
		assertz.StringContains(t, buf.String(), DefaultCompletionSubCommandName)
		assertz.StringContains(t, buf.String(), DefaultGenerateBashCompletionSubCommandName)
	})

	t.Run("failure,open_invalid_file_does_not_exist", func(t *testing.T) {
		c := newTestCommand()
		backup := completionBashTmpl
		t.Cleanup(func() { completionBashTmpl = backup })
		completionBashTmpl = "invalid"
		err := c.Exec(context.Background(), []string{"main-cli", "completion", "bash"})
		requirez.ErrorContains(t, err, `open invalid: file does not exist`)
	})
}
