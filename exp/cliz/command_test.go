package cliz

import (
	"context"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func discard[T any](v T, _ error) T { return v }

func newTestCommand() (cmd *Command, expectedCalledCommands []string, expectedRemainingArgs []string) {
	return &Command{
			Name:        "main-cli",
			Description: "my main command",
			Options: []Option{
				&BoolOption{Name: "bool-opt", Aliases: []string{"b"}, Environment: "BOOL_OPT", Description: "my bool option"},
				&Int64Option{Name: "int64-opt", Aliases: []string{"i64", "int-opt"}, Environment: "INT64_OPT", Description: "my int64 option"},
				&StringOption{Name: "foo", Environment: "FOO", Description: "my foo option"},
			},
			SubCommands: []*Command{
				{
					Name:        "sub-cmd",
					Description: "my sub command",
					Options: []Option{
						&StringOption{Name: "bar", Environment: "BAR", Description: "my bar option"},
					},
					SubCommands: []*Command{
						{
							Name:        "sub-sub-cmd",
							Description: "my sub sub command",
							Options: []Option{
								&StringOption{Name: "baz", Environment: "BAZ", Description: "my baz option"},
							},
						},
					},
				},
				{
					Name:        "sub-cmd2",
					Description: "my sub command2",
				},
				{
					Name:        "sub-cmd3",
					Group:       "groupA",
					Description: "my sub command3",
				},
				{
					Name:        "sub-cmd4",
					Group:       "groupA",
					Description: "my sub command4",
				},
				{
					Name:        "sub-cmd5",
					Group:       "groupB",
					Description: "my sub command5",
				},
				{
					Name:        "sub-cmd6",
					Group:       "groupB",
					Description: "my sub command6",
				},
			},
		},
		[]string{"main-cli", "sub-cmd", "sub-sub-cmd"},
		[]string{"--not-option", "arg1", "arg2"}
}

func TestCommand_parse(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		cmd, expectedCalledCommands, expectedRemainingArgs := newTestCommand()
		actualRemainingArgs, err := cmd.parseArgs(context.Background(), []string{"main-cli", "-b", "--i64", "64", "--foo", "fooValue", "sub-cmd", "--bar", "barValue", "sub-sub-cmd", "--baz", "bazValue", "--", "--not-option", "arg1", "arg2"})
		requirez.NoError(t, err)
		requirez.Equal(t, expectedRemainingArgs, actualRemainingArgs)
		requirez.True(t, discard(cmd.GetOptionBool("bool-opt")))
		requirez.True(t, discard(cmd.GetOptionBool("b")))
		requirez.Equal(t, int64(64), discard(cmd.GetOptionInt64("int-opt")))
		requirez.Equal(t, int64(64), discard(cmd.GetOptionInt64("i64")))
		requirez.Equal(t, "fooValue", discard(cmd.GetOptionString("foo")))
		requirez.Equal(t, "barValue", discard(cmd.GetOptionString("bar")))
		requirez.Equal(t, "bazValue", discard(cmd.GetOptionString("baz")))
		requirez.Equal(t, expectedCalledCommands, cmd.GetCalledCommands())
	})

	t.Run("failure,ErrUnknownOption", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--foo", "fooValue", "--bar", "barValue"})
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})
}
