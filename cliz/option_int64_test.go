package cliz

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hakadoriya/z.go/testingz/assertz"
	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestCommand_GetOptionInt64(t *testing.T) {
	t.Parallel()

	t.Run("error,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		osArgs := []string{"main-cli", "sub-cmd"}
		_, err := c.parse(context.Background(), osArgs)
		requirez.NoError(t, err)
		o, err := c.GetOptionInt64("UNKNOWN")
		assertz.Equal(t, int64(0), o)
		assertz.ErrorIs(t, err, ErrUnknownOption)
		assertz.ErrorContains(t, err, fmt.Sprintf("cmd = %s: option = UNKNOWN: %s", strings.Join(osArgs, " "), ErrUnknownOption))
	})
}

func TestInt64Option_IsZero(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		o := &Int64Option{}
		requirez.True(t, o.IsZero())
	})
}
