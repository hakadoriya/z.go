package mapz

import (
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		source := map[string]int{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		}

		expected := map[string]int{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		}

		actual := Copy(source)

		// copied map should be equal to the source map
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("❌: expected(%+v) != actual(%+v)", expected, actual)
		}

		// source changes should not affect the copied map
		source["foo"] = 100
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("❌: expected(%+v) != actual(%+v)", expected, actual)
		}
	})
}
