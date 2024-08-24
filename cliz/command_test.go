package cliz

import (
	"context"
	"io"
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func discard[T any](v T, _ error) T { return v }

func newTestCommand() (c *Command) {
	return &Command{
		Name:        "main-cli",
		Description: "my main command",
		Options: []Option{
			&StringOption{Name: "string-opt", Aliases: []string{"s"}, Environment: "STRING_OPT", Description: "my string-opt option"},
			&StringOption{Name: "string-opt2"},
			&BoolOption{Name: "bool-opt", Aliases: []string{"b"}, Environment: "BOOL_OPT", Description: "my bool-opt option"},
			&BoolOption{Name: "bool-opt2"},
			&Int64Option{Name: "int64-opt", Aliases: []string{"i64", "int-opt"}, Environment: "INT64_OPT", Description: "my int64-opt option"},
			&Int64Option{Name: "int64-opt2"},
			&Float64Option{Name: "float64-opt", Aliases: []string{"f64", "float-opt"}, Environment: "FLOAT64_OPT", Description: "my float64-opt option"},
			&Float64Option{Name: "float64-opt2"},
			&StringOption{Name: "foo", Environment: "FOO", Description: "my foo option"},
		},
		SubCommands: []*Command{
			{
				Name:        "sub-cmd",
				Aliases:     []string{"sub"},
				Description: "my sub command",
				Options: []Option{
					&StringOption{Name: "bar", Environment: "BAR", Description: "my bar option"},
					&StringOption{Name: "alias", Aliases: []string{"a"}, Description: "my alias option"},
					&StringOption{Name: "string-opt3"},
					&BoolOption{Name: "bool-opt3"},
					&Int64Option{Name: "int64-opt3"},
					&Float64Option{Name: "float64-opt3"},
				},
				ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
					return nil
				},
				SubCommands: []*Command{
					{
						Name:        "sub-sub-cmd",
						Aliases:     []string{"sub-sub"},
						Description: "my sub-sub command",
						Options: []Option{
							&StringOption{Name: "id", Required: true, Description: "my id option"},
							&StringOption{Name: "baz", Environment: "BAZ", Description: "my baz option"},
						},
						PreHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
						ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
						PostHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
					},
					{
						Name:        "sub-sub-cmd2",
						Aliases:     []string{"sub-sub2"},
						Description: "my sub-sub command2",
						Options: []Option{
							&StringOption{Name: "id", Required: true, Description: "my id option"},
						},
						PreHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return io.ErrUnexpectedEOF
						},
						ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
						PostHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
					},
					{
						Name:        "sub-sub-cmd3",
						Aliases:     []string{"sub-sub3"},
						Description: "my sub-sub command3",
						Options: []Option{
							&StringOption{Name: "id", Required: true, Description: "my id option"},
						},
						PreHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
						ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return io.ErrUnexpectedEOF
						},
						PostHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
					},
					{
						Name:        "sub-sub-cmd4",
						Aliases:     []string{"sub-sub4"},
						Description: "my sub-sub command4",
						Options: []Option{
							&StringOption{Name: "id", Required: true, Description: "my id option"},
						},
						PreHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
						ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return nil
						},
						PostHookExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
							return io.ErrUnexpectedEOF
						},
					},
				},
			},
			{
				Name:        "sub-cmd2",
				Description: "my sub command2",
				ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
					return nil
				},
			},
			{
				Name:        "sub-cmd3",
				Group:       "groupA",
				Description: "my sub command3",
				ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
					return nil
				},
			},
			{
				Name:        "sub-cmd4",
				Group:       "groupA",
				Description: "my sub command4",
				ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
					return nil
				},
			},
			{
				Name:        "sub-cmd5",
				Group:       "groupB",
				Description: "my sub command5",
				ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
					return nil
				},
			},
			{
				Name:        "sub-cmd6",
				Group:       "groupB",
				Description: "my sub command6",
				ExecFunc: func(ctx context.Context, rootCmd Cmd, args []string) error {
					return nil
				},
			},
		},
	}
}

func TestCommand_private(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		(*Command)(nil).private()
	})
}

func TestCommand_GetExecutedCommandNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		actual := (*Command)(nil).GetExecutedCommandNames()
		requirez.Equal(t, ([]string)(nil), actual)
	})
}

func TestCommand_GetExecutedCommand(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()

		actual := (*Command)(nil).GetExecutedCommand()
		requirez.Equal(t, (*Command)(nil), actual)
	})

	t.Run("success,empty", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		actual := c.GetExecutedCommand()
		requirez.Equal(t, (*Command)(nil), actual)
	})
}

func TestCommand_Is(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		actual := (*Command)(nil).is("")
		requirez.False(t, actual)
	})
}

func TestCommand_getSubcommand(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		actual := (*Command)(nil).getSubcommand("")
		requirez.Equal(t, (*Command)(nil), actual)
	})
}
