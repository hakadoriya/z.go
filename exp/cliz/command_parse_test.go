package cliz

import (
	"bytes"
	"context"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/assertz"
	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func TestCommand_Parse(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		cmd := newTestCommand()
		expectedCalledCommands, expectedRemainingArgs := []string{"main-cli", "sub-cmd", "sub-sub-cmd"}, []string{"--not-option", "arg1", "arg2"}
		actualRemainingArgs, err := cmd.Parse(context.Background(), []string{
			"main-cli",
			"-s=stringOptValue",
			"--string-opt2=string-opt2-value",
			"-b",
			"--bool-opt2=true",
			"--i64", "64",
			"--int64-opt2=-128",
			"--float64-opt", "3.14",
			"--float64-opt2=6.28",
			"--foo", "fooValue",
			"sub-cmd",
			"--bar", "barValue",
			"sub-sub",
			"--id=idValue",
			"--baz", "bazValue",
			"--",
			"--not-option",
			"arg1",
			"arg2",
		})
		requirez.NoError(t, err)
		assertz.Equal(t, expectedRemainingArgs, actualRemainingArgs)
		assertz.True(t, discard(cmd.GetOptionBool("bool-opt")))
		assertz.True(t, discard(cmd.GetOptionBool("b")))
		assertz.True(t, discard(cmd.GetOptionBool("bool-opt2")))
		assertz.Equal(t, int64(64), discard(cmd.GetOptionInt64("int-opt")))
		assertz.Equal(t, int64(64), discard(cmd.GetOptionInt64("i64")))
		assertz.Equal(t, int64(-128), discard(cmd.GetOptionInt64("int64-opt2")))
		assertz.Equal(t, "stringOptValue", discard(cmd.GetOptionString("string-opt")))
		assertz.Equal(t, "stringOptValue", discard(cmd.GetOptionString("s")))
		assertz.Equal(t, "string-opt2-value", discard(cmd.GetOptionString("string-opt2")))
		assertz.Equal(t, "fooValue", discard(cmd.GetOptionString("foo")))
		assertz.Equal(t, "barValue", discard(cmd.GetOptionString("bar")))
		assertz.Equal(t, "idValue", discard(cmd.GetOptionString("id")))
		assertz.Equal(t, "bazValue", discard(cmd.GetOptionString("baz")))
		assertz.Equal(t, expectedCalledCommands, cmd.GetCalledCommands())
	})

	t.Run("failure,StringOption,ErrMissingOptionValue", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--foo"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("failure,BoolOption,strconv.ParseBool", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--bool-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Int64Option,ErrMissingOptionValue", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--int64-opt"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("failure,Int64Option,argIsHyphenOption,strconv.ParseInt", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--int64-opt", "FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Int64Option,argIsHyphenOptionEqual,strconv.ParseInt", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--int64-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Float64Option,ErrMissingOptionValue", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--float64-opt"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("failure,Float64Option,argIsHyphenOption,strconv.ParseFloat", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--float64-opt", "FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,Float64Option,argIsHyphenOptionEqual,strconv.ParseFloat", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--float64-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,HelpOption,strconv.ParseBool", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--help=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "FAILURE": invalid syntax`)
	})

	t.Run("failure,ErrInvalidOptionType", func(t *testing.T) {
		cmd := newTestCommand()
		cmd.Options = append(cmd.Options, &unknownOptionType{Name: "unknown", Environment: "UNKNOWN", Description: "my unknown option"})
		_, err := cmd.Parse(context.Background(), []string{"main-cli", "--foo", "fooValue", "--unknown", "unknownValue"})
		requirez.ErrorIs(t, err, ErrInvalidOptionType)
	})
}

func TestCommand_parse(t *testing.T) {
	t.Parallel()

	t.Run("success,HelpOption", func(t *testing.T) {
		cmd := newTestCommand()
		buf := new(bytes.Buffer)
		cmd.Stderr = buf
		_, err := cmd.parse(context.Background(), []string{"main-cli", "--help=true"})
		requirez.ErrorIs(t, err, ErrHelp)
	})
}

func TestCommand_parseArgs(t *testing.T) {
	t.Parallel()

	t.Run("failure,ErrInvalidOptionType", func(t *testing.T) {
		cmd := newTestCommand()
		cmd.Options = append(cmd.Options, &unknownOptionType{Name: "unknown", Environment: "UNKNOWN", Description: "my unknown option"})
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--foo", "fooValue", "--unknown", "unknownValue"})
		requirez.ErrorIs(t, err, ErrInvalidOptionType)
	})

	t.Run("failure,ErrUnknownOption", func(t *testing.T) {
		cmd := newTestCommand()
		_, err := cmd.parseArgs(context.Background(), []string{"main-cli", "--foo", "fooValue", "sub-cmd", "--unknown", "unknownValue"})
		t.Logf("err: %v", err)
		requirez.ErrorIs(t, err, ErrUnknownOption)
	})
}
