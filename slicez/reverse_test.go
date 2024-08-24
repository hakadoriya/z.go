package slicez

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		source := []int{1, 2, 3, 4, 5}
		expected := []int{5, 4, 3, 2, 1}
		actual := Reverse(source)
		if len(expected) != len(actual) {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})
}
