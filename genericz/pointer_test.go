package genericz

import (
	"testing"
)

func TestPointer(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		const expected = "test"
		actual := Pointer(expected)
		if expected != *actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, *actual)
		}
	})
}

func TestPtr(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		const expected = "test"
		actual := Ptr(expected)
		if expected != *actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, *actual)
		}
	})
}
