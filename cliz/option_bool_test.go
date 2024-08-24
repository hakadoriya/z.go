package cliz

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hakadoriya/z.go/cliz/clierrz"
	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_GetOptionBool(t *testing.T) {
	t.Parallel()

	t.Run("error,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "sub-cmd"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		o, err := c.GetOptionBool("UNKNOWN")
		assertz.False(t, o)
		assertz.ErrorIs(t, err, clierrz.ErrUnknownOption)
		assertz.ErrorContains(t, err, fmt.Sprintf("cmd = %s: option = UNKNOWN: %s", strings.Join(osArgs, " "), clierrz.ErrUnknownOption))
	})
}

func TestBoolOption_IsZero(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		o := &BoolOption{}
		requirez.True(t, o.IsZero())
	})
}
