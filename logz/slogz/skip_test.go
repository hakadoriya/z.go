package slogz

import (
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/requirez"
)

func Test_addCallerSkip(t *testing.T) {
	t.Parallel()

	t.Run("failure,0", func(t *testing.T) {
		t.Parallel()

		actual := addCallerSkip(nil)
		requirez.Equal(t, 0, actual)
	})
}
