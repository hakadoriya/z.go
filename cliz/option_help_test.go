package cliz

import (
	"testing"

	"github.com/hakadoriya/z.go/testingz/requirez"
)

func TestHelpOption_IsZero(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		o := &HelpOption{}
		requirez.True(t, o.IsZero())
	})
}

func TestHelpOption_GetDescription(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		const expected = "EXPECTED"
		o := &HelpOption{
			Description: expected,
		}
		requirez.Equal(t, expected, o.GetDescription())
	})
}
