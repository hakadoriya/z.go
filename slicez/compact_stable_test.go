package slicez

import (
	"reflect"
	"testing"
)

func TestCompactStable(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		src := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5}
		expect := []int{1, 2, 3, 4, 5}
		actual := CompactStable(src)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})
}
