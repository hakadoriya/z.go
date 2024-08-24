package cliz

import (
	"bytes"
	"context"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func TestCommand_checkHelp(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cmd, _, _ := newTestCommand()

		buf := new(bytes.Buffer)

		ctx := WithContext(context.Background(), cmd)
		ctx = stdoutWithContext(ctx, buf)
		ctx = stderrWithContext(ctx, buf)
		actualRemainingArgs, err := cmd.parse(ctx, []string{"main-cli", "--help"})
		requirez.ErrorIs(t, err, ErrHelp)
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
        my string-opt2 option
    --bool-opt, -b (env: BOOL_OPT, default: false)
        my bool-opt option
    --bool-opt2 (default: false)
        my bool-opt2 option
    --int64-opt, -i64, -int-opt (env: INT64_OPT, default: 0)
        my int64-opt option
    --int64-opt2 (default: 0)
        my int64-opt2 option
    --float64-opt, -f64, -float-opt (env: FLOAT64_OPT, default: 0)
        my float64-opt option
    --float64-opt2 (default: 0)
        my float64-opt2 option
    --foo (env: FOO, default: )
        my foo option
    --help (default: false)
        show usage
`
		requirez.Equal(t, usage, buf.String())
	})
}
