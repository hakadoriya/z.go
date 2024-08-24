package stringz

import "testing"

func TestTrimLeftSpace(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		const expect = "a"
		actual := TrimLeftSpace("   	a")
		if expect != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expect, actual)
		}
	})
}

func TestTrimRightSpace(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		const expect = "a"
		actual := TrimRightSpace("a	   ")
		if expect != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expect, actual)
		}
	})
}
