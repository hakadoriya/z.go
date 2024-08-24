package cliz

import (
	"context"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

//nolint:paralleltest
func TestCommand_loadEnvironments(t *testing.T) {
	t.Run("success,", func(t *testing.T) {
		c := newTestCommand()
		err := c.loadEnvironments()
		requirez.NoError(t, err)
	})

	t.Run("error,ErrInvalidOptionType", func(t *testing.T) {
		c := newTestCommand()
		for _, subcmd := range c.SubCommands {
			subcmd.Options = append(subcmd.Options, &unknownOptionType{Name: "unknown", Environment: "UNKNOWN", Description: "my unknown option"})
		}
		err := c.loadEnvironments()
		requirez.ErrorIs(t, err, ErrInvalidOptionType)
	})

	t.Run("success,StringOption", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("STRING_OPT", "envValue")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.NoError(t, err)
		assertz.Equal(t, "envValue", discard(c.GetOptionString("string-opt")))
	})

	t.Run("success,BoolOption", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("BOOL_OPT", "TRUE")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.NoError(t, err)
		assertz.Equal(t, true, discard(c.GetOptionBool("bool-opt")))
	})

	t.Run("error,BoolOption", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("BOOL_OPT", "FAILURE")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorContains(t, err, `strconv.ParseBool: parsing "FAILURE": invalid syntax`)
	})

	t.Run("success,Int64Option", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("INT64_OPT", "64")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.NoError(t, err)
		assertz.Equal(t, int64(64), discard(c.GetOptionInt64("int64-opt")))
	})

	t.Run("error,Int64Option", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("INT64_OPT", "FAILURE")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorContains(t, err, `strconv.ParseInt: parsing "FAILURE": invalid syntax`)
	})

	t.Run("success,Uint64Option", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("UINT64_OPT", "64")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.NoError(t, err)
		assertz.Equal(t, uint64(64), discard(c.GetOptionUint64("uint64-opt")))
	})

	t.Run("error,Uint64Option", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("UINT64_OPT", "FAILURE")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorContains(t, err, `strconv.ParseUint: parsing "FAILURE": invalid syntax`)
	})

	t.Run("success,Float64Option", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("FLOAT64_OPT", "3.1416926535")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.NoError(t, err)
		assertz.Equal(t, 3.1416926535, discard(c.GetOptionFloat64("float64-opt")))
	})

	t.Run("error,Float64Option", func(t *testing.T) {
		c := newTestCommand()
		t.Setenv("FLOAT64_OPT", "FAILURE")
		_, err := c.parse(context.Background(), []string{"main-cli"})
		requirez.ErrorContains(t, err, `strconv.ParseFloat: parsing "FAILURE": invalid syntax`)
	})
}
