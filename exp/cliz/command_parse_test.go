package cliz

import (
	"context"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

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

	t.Run("failure,StringOption,ErrMissingOptionValue", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--foo"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("failure,ErrUnknownOptiona", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		cmd.Options = append(cmd.Options, &unknownOptionType{Name: "unknown", Environment: "UNKNOWN", Description: "my unknown option"})
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--foo", "fooValue", "--unknown", "unknownValue"})
		requirez.ErrorIs(t, err, ErrInvalidOptionType)
	})

	t.Run("failure,ErrUnknownOption", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--foo", "fooValue", "sub-cmd", "--unknown", "unknownValue"})
		t.Logf("err: %v", err)
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})
}
