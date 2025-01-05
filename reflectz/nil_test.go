package reflectz

import (
	"net/http"
	"testing"
)

func TestIsNil(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()

		const expected = true
		actual := IsNil(nil)
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})

	t.Run("success,[]byte", func(t *testing.T) {
		t.Parallel()

		const expected = false
		actual := IsNil([]byte{})
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})

	t.Run("success,http.Request", func(t *testing.T) {
		t.Parallel()

		const expected = true
		actual := IsNil((*http.Request)(nil))
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})

	t.Run("success,http.ResponseWriter", func(t *testing.T) {
		t.Parallel()

		const expected = true
		actual := IsNil((*http.ResponseWriter)(nil))
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})
}
