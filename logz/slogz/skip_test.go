package slogz

import (
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func Test_addCallerSkip(t *testing.T) {
	t.Parallel()

	t.Run("error,0", func(t *testing.T) {
		t.Parallel()

		actual := addCallerSkip(nil)
		requirez.Equal(t, 0, actual)
	})
}
