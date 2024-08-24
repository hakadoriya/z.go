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
		args := []string{
			"main-cli",
			"-s", "stringOptValue",
			"--string-opt2=string-opt2-value",
			"-b",
			"--bool-opt2=true",
			"--i64", "64",
			"--int64-opt2=128",
			"--float64-opt", "3.14",
			"--float64-opt2=6.28",
			"--foo", "fooValue",
			"sub-cmd",
			"--bar", "barValue",
			"sub-sub-cmd",
			"--baz", "bazValue",
			"--",
			"--not-option",
			"arg1",
			"arg2",
		}
		actualRemainingArgs, err := cmd.parseArgs(context.Background(), args)
		requirez.NoError(t, err)
		requirez.Equal(t, expectedRemainingArgs, actualRemainingArgs)
		requirez.True(t, discard(cmd.GetOptionBool("bool-opt")))
		requirez.True(t, discard(cmd.GetOptionBool("b")))
		requirez.True(t, discard(cmd.GetOptionBool("bool-opt2")))
		requirez.Equal(t, int64(64), discard(cmd.GetOptionInt64("int-opt")))
		requirez.Equal(t, int64(64), discard(cmd.GetOptionInt64("i64")))
		requirez.Equal(t, int64(128), discard(cmd.GetOptionInt64("int64-opt2")))
		requirez.Equal(t, "stringOptValue", discard(cmd.GetOptionString("string-opt")))
		requirez.Equal(t, "stringOptValue", discard(cmd.GetOptionString("s")))
		requirez.Equal(t, "string-opt2-value", discard(cmd.GetOptionString("string-opt2")))
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

	t.Run("failure,BoolOption,strconv.ParseBool", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--bool-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Int64Option,ErrMissingOptionValue", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--int64-opt"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("failure,Int64Option,argIsHyphenOption,strconv.ParseInt", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--int64-opt", "FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Int64Option,argIsHyphenOptionEqual,strconv.ParseInt", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--int64-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Float64Option,ErrMissingOptionValue", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--float64-opt"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("failure,Float64Option,argIsHyphenOption,strconv.ParseFloat", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--float64-opt", "FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Float64Option,argIsHyphenOptionEqual,strconv.ParseFloat", func(t *testing.T) {
		cmd, _, _ := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--float64-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,ErrInvalidOptionType", func(t *testing.T) {
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
