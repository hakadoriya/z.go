package panicz_test

import (
	"io"
	"testing"

	"github.com/hakadoriya/z.go/panicz"
)

func TestPanic(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("❌: %s: expected not to panic", t.Name())
			}
		}()

		panicz.Panic(nil, panicz.WithPanicOptionIgnoreErrors(io.EOF))
		panicz.Panic(io.EOF, panicz.WithPanicOptionIgnoreErrors(io.EOF))
	})

	t.Run("failure,", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("❌: %s: expected to panic", t.Name())
			}
		}()

		panicz.Panic(io.EOF)
	})
}
