package cliz

import (
	"bytes"
	"os"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func TestCommand_GetStdout(t *testing.T) {
	t.Parallel()

	t.Run("success,stdout", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		buf := new(bytes.Buffer)
		c.SetStdout(buf)
		actual := c.GetStdout()
		requirez.Equal(t, buf, actual)
	})

	t.Run("success,os.Stdout", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		actual := c.GetStdout()
		requirez.Equal(t, os.Stdout, actual)
	})
}

func TestCommand_GetStderr(t *testing.T) {
	t.Parallel()

	t.Run("success,stderr", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		buf := new(bytes.Buffer)
		c.SetStderr(buf)
		actual := c.GetStderr()
		requirez.Equal(t, buf, actual)
	})

	t.Run("success,os.Stderr", func(t *testing.T) {
		t.Parallel()

		c := newTestCommand()
		actual := c.GetStderr()
		requirez.Equal(t, os.Stderr, actual)
	})
}

func TestCommand_SetStdout(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		(*Command)(nil).SetStdout(nil)
	})
}

func TestCommand_SetStdoutRecursive(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		(*Command)(nil).SetStdoutRecursive(nil)
	})
}

func TestCommand_SetStderr(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		(*Command)(nil).SetStderr(nil)
	})
}

func TestCommand_SetStderrRecursive(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		(*Command)(nil).SetStderrRecursive(nil)
	})
}
