package slicez

import (
	"reflect"
	"regexp"
	"testing"
)

func TestDistinct(t *testing.T) {
	t.Parallel()

	t.Run("success,case1", func(t *testing.T) {
		t.Parallel()

		expect := []int{1, 2, 3}
		actual := Distinct([]int{1, 2, 2, 3, 3, 3})
		if len(expect) != len(actual) {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})
}

func TestDeepDistinct(t *testing.T) {
	t.Parallel()

	t.Run("success,case1", func(t *testing.T) {
		t.Parallel()

		expect := []interface{}{1, 2, 3, regexp.MustCompile(".*")}
		actual := DeepDistinct([]interface{}{1, 2, 2, 3, 3, 3, regexp.MustCompile(".*"), regexp.MustCompile(".*")})
		if len(expect) != len(actual) {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})
}
