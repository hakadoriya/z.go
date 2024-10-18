package genericz

import (
	"testing"
)

func TestZero(t *testing.T) {
	t.Parallel()

	t.Run("success,string", func(t *testing.T) {
		t.Parallel()
		const expected = ""
		actual := Zero("1")
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})
	t.Run("success,int", func(t *testing.T) {
		t.Parallel()
		const expected = 0
		actual := Zero(1)
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})
}

func TestSliceElemZero(t *testing.T) {
	t.Parallel()

	t.Run("success,string", func(t *testing.T) {
		t.Parallel()
		const expected = ""
		actual := SliceElemZero([]string{"1"})
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("success,int", func(t *testing.T) {
		t.Parallel()
		const expected = 0
		actual := SliceElemZero([]int{1})
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	t.Run("success,string,true", func(t *testing.T) {
		t.Parallel()
		const expected = true
		actual := IsZero("")
		if expected != actual {
			t.Errorf("❌: expected(%t) != actual(%t)", expected, actual)
		}
	})

	t.Run("success,string,false", func(t *testing.T) {
		t.Parallel()
		const expected = false
		actual := IsZero("1")
		if expected != actual {
			t.Errorf("❌: expected(%t) != actual(%t)", expected, actual)
		}
	})

	t.Run("success,int,true", func(t *testing.T) {
		t.Parallel()
		const expected = true
		actual := IsZero(0)
		if expected != actual {
			t.Errorf("❌: expected(%t) != actual(%t)", expected, actual)
		}
	})

	t.Run("success,int,false", func(t *testing.T) {
		t.Parallel()
		const expected = false
		actual := IsZero(1)
		if expected != actual {
			t.Errorf("❌: expected(%t) != actual(%t)", expected, actual)
		}
	})
}
