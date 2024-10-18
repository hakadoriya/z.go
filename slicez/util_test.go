package slicez

import (
	"reflect"
	"testing"
)

func TestFirst(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "1"
		s := []string{"1", "2", "3"}
		actual := First(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("not", func(t *testing.T) {
		t.Parallel()
		expect := ""
		s := []string{}
		actual := First(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestLast(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "3"
		s := []string{"1", "2", "3"}
		actual := Last(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("not", func(t *testing.T) {
		t.Parallel()
		expect := ""
		s := []string{}
		actual := Last(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}
