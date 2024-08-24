package cliz

import (
	"bytes"
	"context"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_Parse(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{
			"main-cli",
			"-s=stringOptValue",
			"--string-opt2=string-opt2-value",
			"-b",
			"--bool-opt2=true",
			"--i64", "64",
			"--int64-opt2=-128",
			"-u64", "256",
			"--uint64-opt2=512",
			"--f64", "3.14",
			"--float64-opt2=6.28",
			"--hidden-opt=hiddenOptValue",
			"--foo", "fooValue",
			"sub-cmd",
			"--bar", "barValue",
			"--string-opt3", "stringOpt3Value",
			"--bool-opt3",
			"--int64-opt3", "128",
			"--uint64-opt3", "256",
			"--float64-opt3", "9.42",
			"sub-sub",
			"--id=idValue",
			"--baz", "bazValue",
			"--",
			"--not-option",
			"arg1",
			"arg2",
		}
		expectedCalledCommands, expectedRemainingArgs := []string{"main-cli", "sub-cmd", "sub-sub-cmd"}, []string{"--not-option", "arg1", "arg2"}
		actualRemainingArgs, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		assertz.Equal(t, expectedRemainingArgs, actualRemainingArgs)
		assertz.Equal(t, "stringOptValue", discard(c.GetOptionString("string-opt")))
		assertz.Equal(t, "stringOptValue", discard(c.GetOptionString("s")))
		assertz.Equal(t, "string-opt2-value", discard(c.GetOptionString("string-opt2")))
		assertz.True(t, discard(c.GetOptionBool("bool-opt")))
		assertz.True(t, discard(c.GetOptionBool("b")))
		assertz.True(t, discard(c.GetOptionBool("bool-opt2")))
		assertz.Equal(t, int64(64), discard(c.GetOptionInt64("int-opt")))
		assertz.Equal(t, int64(64), discard(c.GetOptionInt64("i64")))
		assertz.Equal(t, int64(-128), discard(c.GetOptionInt64("int64-opt2")))
		assertz.Equal(t, uint64(256), discard(c.GetOptionUint64("uint-opt")))
		assertz.Equal(t, uint64(256), discard(c.GetOptionUint64("u64")))
		assertz.Equal(t, uint64(512), discard(c.GetOptionUint64("uint64-opt2")))
		assertz.Equal(t, 3.14, discard(c.GetOptionFloat64("float64-opt")))
		assertz.Equal(t, 3.14, discard(c.GetOptionFloat64("f64")))
		assertz.Equal(t, 6.28, discard(c.GetOptionFloat64("float64-opt2")))
		assertz.Equal(t, "hiddenOptValue", discard(c.GetOptionString("hidden-opt")))
		assertz.Equal(t, "fooValue", discard(c.GetOptionString("foo")))
		assertz.Equal(t, "barValue", discard(c.GetOptionString("bar")))
		assertz.Equal(t, "stringOpt3Value", discard(c.GetOptionString("string-opt3")))
		assertz.True(t, discard(c.GetOptionBool("bool-opt3")))
		assertz.Equal(t, int64(128), discard(c.GetOptionInt64("int64-opt3")))
		assertz.Equal(t, uint64(256), discard(c.GetOptionUint64("uint64-opt3")))
		assertz.Equal(t, 9.42, discard(c.GetOptionFloat64("float64-opt3")))
		assertz.Equal(t, "idValue", discard(c.GetOptionString("id")))
		assertz.Equal(t, "bazValue", discard(c.GetOptionString("baz")))
		assertz.Equal(t, expectedCalledCommands, c.GetExecutedCommandNames())

		actualRemainingArgs2, err2 := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err2)
		assertz.Equal(t, expectedRemainingArgs, actualRemainingArgs2)
	})

	t.Run("success,HelpOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		buf := new(bytes.Buffer)
		c.SetStdoutRecursive(buf)
		c.SetStderrRecursive(buf)
		_, err := c.parse(context.Background(), []string{"main-cli", "--help=true"})
		requirez.ErrorIs(t, err, ErrHelp)
	})

	t.Run("error,StringOption,ErrMissingOptionValue", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--foo"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("error,preCheckSubCommands,ErrDuplicateSubCommand", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		c.SubCommands = append(c.SubCommands, &Command{Name: "sub-cmd"})
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorIs(t, err, ErrDuplicateSubCommand)
	})

	t.Run("error,preCheckOptions,ErrDuplicateOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		c.Options = append(c.Options, &StringOption{Name: "foo"})
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorIs(t, err, ErrDuplicateOption)
	})

	t.Run("error,postCheckOptions,ErrOptionRequired", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		c.Options = append(c.Options, &StringOption{Name: "require", Required: true, Description: "my require option"})
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorIs(t, err, ErrOptionRequired)
	})

	t.Run("error,BoolOption,strconv.ParseBool", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--bool-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,Int64Option,ErrMissingOptionValue", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--int64-opt"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("error,Int64Option,argIsHyphenOption,strconv.ParseInt", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--int64-opt", "FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,Int64Option,argIsHyphenOptionEqual,strconv.ParseInt", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--int64-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,Uint64Option,ErrMissingOptionValue", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--uint64-opt"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("error,Uint64Option,argIsHyphenOption,strconv.ParseUint", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--uint64-opt", "FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseUint: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,Uint64Option,argIsHyphenOptionEqual,strconv.ParseUint", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--uint64-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseUint: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,Float64Option,ErrMissingOptionValue", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--float64-opt"})
		requirez.ErrorIs(t, err, ErrMissingOptionValue)
	})

	t.Run("error,Float64Option,argIsHyphenOption,strconv.ParseFloat", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--float64-opt", "FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,Float64Option,argIsHyphenOptionEqual,strconv.ParseFloat", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--float64-opt=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,HelpOption,strconv.ParseBool", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parse(context.Background(), []string{"main-cli", "--help=FAILURE"})
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "FAILURE": invalid syntax`)
	})

	t.Run("error,ErrInvalidOptionType", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		c.Options = append(c.Options, &unknownOptionType{Name: "unknown", Environment: "UNKNOWN", Description: "my unknown option"})
		_, err := c.parse(context.Background(), []string{"main-cli", "--foo", "fooValue", "--unknown", "unknownValue"})
		requirez.ErrorIs(t, err, ErrInvalidOptionType)
	})

	t.Run("error,context.Canceled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		c := newTestCommand()
		_, err := c.parse(ctx, []string{"main-cli"})
		requirez.ErrorIs(t, err, context.Canceled)
	})
}

func TestCommand_Parse_Environment(t *testing.T) {
	t.Run("error,strconv.ParseBool", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("BOOL_OPT", "FAILURE")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "FAILURE": invalid syntax`)
	})
}

func TestCommand_parseArgs(t *testing.T) {
	t.Parallel()

	t.Run("error,ErrInvalidOptionType", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		c.Options = append(c.Options, &unknownOptionType{Name: "unknown", Environment: "UNKNOWN", Description: "my unknown option"})
		_, err := c.parseArgs([]string{"main-cli", "--foo", "fooValue", "--unknown", "unknownValue"})
		requirez.ErrorIs(t, err, ErrInvalidOptionType)
	})

	t.Run("error,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		_, err := c.parseArgs([]string{"main-cli", "--foo", "fooValue", "sub-cmd", "--unknown", "unknownValue"})
		requirez.ErrorIs(t, err, ErrUnknownOption)
		assertz.ErrorContains(t, err, "main-cli: sub-cmd: --unknown: unknown option")
	})
}
