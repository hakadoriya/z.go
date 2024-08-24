package cliz

import (
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_checkOptions(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{Name: "foo"},
					},
				},
			},
		}

		err := c.preCheckOptions()
		requirez.NoError(t, err)
	})

	t.Run("error,ErrDuplicateOption,Name", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{Name: "foo"},
						&StringOption{Name: "foo"},
					},
				},
			},
		}

		err := c.preCheckOptions()
		requirez.ErrorIs(t, err, ErrDuplicateOption)
	})

	t.Run("error,ErrDuplicateOption,Aliases", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{Aliases: []string{"f"}},
						&StringOption{Aliases: []string{"f"}},
					},
				},
			},
		}

		err := c.preCheckOptions()
		requirez.ErrorIs(t, err, ErrDuplicateOption)
	})

	t.Run("error,ErrDuplicateOption,Environment", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{Env: "FOO"},
						&StringOption{Env: "FOO"},
					},
				},
			},
		}

		err := c.preCheckOptions()
		requirez.ErrorIs(t, err, ErrDuplicateOption)
	})
}
