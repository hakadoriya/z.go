package cliz

import (
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func discard[T any](v T, _ error) T { return v }

func newTestCommand() (cmd *Command, expectedCalledCommands []string, expectedRemainingArgs []string) {
	return &Command{
			Name:        "main-cli",
			Description: "my main command",
			Options: []Option{
				&StringOption{Name: "string-opt", Aliases: []string{"s"}, Environment: "STRING_OPT", Description: "my string-opt option"},
				&StringOption{Name: "string-opt2", Description: "my string-opt2 option"},
				&BoolOption{Name: "bool-opt", Aliases: []string{"b"}, Environment: "BOOL_OPT", Description: "my bool-opt option"},
				&BoolOption{Name: "bool-opt2", Description: "my bool-opt2 option"},
				&Int64Option{Name: "int64-opt", Aliases: []string{"i64", "int-opt"}, Environment: "INT64_OPT", Description: "my int64-opt option"},
				&Int64Option{Name: "int64-opt2", Description: "my int64-opt2 option"},
				&Float64Option{Name: "float64-opt", Aliases: []string{"f64", "float-opt"}, Environment: "FLOAT64_OPT", Description: "my float64-opt option"},
				&Float64Option{Name: "float64-opt2", Description: "my float64-opt2 option"},
				&StringOption{Name: "foo", Environment: "FOO", Description: "my foo option"},
				&StringOption{Name: "id", Description: "my id option"},
			},
			SubCommands: []*Command{
				{
					Name:        "sub-cmd",
					Aliases:     []string{"sub"},
					Description: "my sub command",
					Options: []Option{
						&StringOption{Name: "bar", Environment: "BAR", Description: "my bar option"},
						&StringOption{Name: "id", Description: "my id option"},
					},
					SubCommands: []*Command{
						{
							Name:        "sub-sub-cmd",
							Aliases:     []string{"sub-sub"},
							Description: "my sub sub command",
							Options: []Option{
								&StringOption{Name: "baz", Environment: "BAZ", Description: "my baz option"},
								&StringOption{Name: "id", Description: "my id option"},
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

func TestCommand_GetCalledCommands(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		actual := (*Command)(nil).GetCalledCommands()
		requirez.Equal(t, ([]string)(nil), actual)
	})
}

func TestCommand_Is(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		actual := (*Command)(nil).Is("")
		requirez.False(t, actual)
	})
}

func TestCommand_getSubcommand(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		actual := (*Command)(nil).getSubcommand("")
		requirez.Equal(t, nil, actual)
		requirez.Equal(t, (*Command)(nil), actual)
	})
}
