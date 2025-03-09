package slicez

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expect := []int{1, 4, 9}
		actual := Map([]int{1, 2, 3}, func(index int, elem int) int {
			return elem * elem
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})
}
