package slicez_test

import (
	"reflect"
	"testing"

	"github.com/hakadoriya/z.go/slicez"
)

func TestForEach(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2}
		s := []int{1, 2, 3}
		actual := make([]int, 0)
		slicez.ForEach(s, func(index int, elem int) {
			actual = append(actual, elem-1)
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}
