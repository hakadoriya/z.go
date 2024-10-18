package slicez

import (
	"math"
	"testing"
)

func TestContains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := true
		s := []int{0, 1, 2, 3}
		value := 1
		actual := Contains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("not", func(t *testing.T) {
		t.Parallel()
		expect := false
		s := []int{0, 1, 2, 3}
		value := math.MaxInt
		actual := Contains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestDeepContains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := true
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{1}
		actual := DeepContains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("not", func(t *testing.T) {
		t.Parallel()
		expect := false
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{}
		actual := DeepContains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}
