package slicez

import (
	"math"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	t.Parallel()

	t.Run("success,case1", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}
		s := []string{"1", "2", "3", "4", "5", "6"}
		actual := Split(s, 2)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success,case2", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2"}, {"3", "4"}, {"5"}}
		s := []string{"1", "2", "3", "4", "5"}
		actual := Split(s, 2)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success,case3", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2", "3", "4", "5"}}
		s := []string{"1", "2", "3", "4", "5"}
		actual := Split(s, math.MaxInt)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}
