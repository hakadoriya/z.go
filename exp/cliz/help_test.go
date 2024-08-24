package cliz

import (
	"bytes"
	"context"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func TestCommand_checkHelp(t *testing.T) {
	t.Parallel()

	t.Run("success,Command", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		buf := new(bytes.Buffer)

		ctx := WithContext(context.Background(), cmd)
		cmd.SetStdoutRecursive(buf)
		cmd.SetStderrRecursive(buf)
		actualRemainingArgs, err := cmd.parse(ctx, []string{"main-cli", "--help"})
		requirez.True(t, IsHelp(err))
		requirez.Equal(t, 0, len(actualRemainingArgs))
		const usage = `Usage:
    main-cli [options] <subcommand>

Description:
    my main command

Sub commands:
    sub-cmd     my sub command
    sub-cmd2    my sub command2
    groupA:
        sub-cmd3    my sub command3
        sub-cmd4    my sub command4
    groupB:
        sub-cmd5    my sub command5
        sub-cmd6    my sub command6

Options:
    --string-opt, -s (env: STRING_OPT, default: )
        my string-opt option
    --string-opt2 (default: )
        string value of string-opt2
    --bool-opt, -b (env: BOOL_OPT, default: false)
        my bool-opt option
    --bool-opt2 (default: false)
        bool value of bool-opt2
    --int64-opt, -i64, -int-opt (env: INT64_OPT, default: 0)
        my int64-opt option
    --int64-opt2 (default: 0)
        int64 value of int64-opt2
    --float64-opt, -f64, -float-opt (env: FLOAT64_OPT, default: 0)
        my float64-opt option
    --float64-opt2 (default: 0)
        float64 value of float64-opt2
    --foo (env: FOO, default: )
        my foo option
    --help (default: false)
        show usage
`
		requirez.Equal(t, usage, buf.String())
	})

	t.Run("success,SubCommand", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		buf := new(bytes.Buffer)

		ctx := WithContext(context.Background(), cmd)
		cmd.SetStdoutRecursive(buf)
		cmd.SetStderrRecursive(buf)
		actualRemainingArgs, err := cmd.parse(ctx, []string{"main-cli", "sub-cmd", "--help"})
		requirez.True(t, IsHelp(err))
		requirez.Equal(t, 0, len(actualRemainingArgs))
		const usage = `Usage:
    main-cli sub-cmd [options] <subcommand>

Description:
    my sub command

Sub commands:
    sub-sub-cmd     my sub sub command
    sub-sub-cmd2    my sub sub command2

Options:
    --bar (env: BAR, default: )
        my bar option
    --string-opt3 (default: )
        string value of string-opt3
    --bool-opt3 (default: false)
        bool value of bool-opt3
    --int64-opt3 (default: 0)
        int64 value of int64-opt3
    --float64-opt3 (default: 0)
        float64 value of float64-opt3
    --help (default: false)
        show usage
`
		requirez.Equal(t, usage, buf.String())
	})
}

func TestCommand_getHelpOption(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		c.initAppendHelpOption()
		o, ok := c.getHelpOption()
		requirez.True(t, ok)
		requirez.NotNil(t, o)
	})
}

func TestCommand_ShowUsage(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		buf := new(bytes.Buffer)
		c.SetStderr(buf)
		c.UsageFunc = func(cmd *Command) {
			cmd.GetStderr().Write([]byte("UsageFunc\n"))
		}
		c.ShowUsage()
		requirez.Equal(t, "UsageFunc\n", buf.String())
	})
}

func TestCommand_showUsage(t *testing.T) {
	t.Parallel()

	t.Run("success,Usage", func(t *testing.T) {
		t.Parallel()

		c := &Command{Name: "main-cli", Usage: "main-cli [OPTIONS]"}
		buf := new(bytes.Buffer)
		c.SetStderr(buf)
		c.ShowUsage()
		const expected = `Usage:
    main-cli [OPTIONS]
`
		requirez.Equal(t, expected, buf.String())
	})

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		c := &Command{Name: "main-cli", Options: []Option{&StringOption{Name: "required", Required: true}}}
		buf := new(bytes.Buffer)
		c.SetStderr(buf)
		_, err := c.parse(context.Background(), []string{"main-cli", "--help"})
		requirez.ErrorIs(t, err, ErrHelp)
		const expected = `Usage:
    main-cli [options]

Options:
    --required (required, default: )
        string value of required
    --help (default: false)
        show usage
`
		requirez.Equal(t, expected, buf.String())
	})
}
