package cliz

import (
	"testing"

	"github.com/hakadoriya/z.go/cliz/clicorez"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_preCheckSubCommands(t *testing.T) {
	t.Parallel()
	t.Run("error,ErrDuplicateSubCommand", func(t *testing.T) {
		t.Parallel()

		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					SubCommands: []*Command{
						{
							Name: "sub-sub-cmd",
						},
						{
							Name: "sub-sub-cmd",
						},
					},
				},
			},
		}

		err := c.preCheckSubCommands()
		requirez.ErrorIs(t, err, clicorez.ErrDuplicateSubCommand)
	})
}
