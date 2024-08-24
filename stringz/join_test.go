package stringz

import (
	"fmt"
	"testing"
)

func TestJoin(t *testing.T) {
	t.Parallel()

	t.Run("success,normal", func(t *testing.T) {
		t.Parallel()

		const expect = "a_ _b_ _c"
		actual := Join("_ _", "a", "b", "c")

		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})
}

var _ fmt.Stringer = (*Stringer)(nil)

type Stringer struct {
	StringFunc func() string
}

func (r *Stringer) String() string {
	if r.StringFunc == nil {
		return ""
	}
	return r.StringFunc()
}

func TestJoinStringers(t *testing.T) {
	t.Parallel()

	t.Run("success,normal", func(t *testing.T) {
		t.Parallel()

		const expect = "a_ _b_ _c"
		actual := JoinStringers(
			"_ _",
			&Stringer{StringFunc: func() string { return "a" }},
			&Stringer{StringFunc: func() string { return "b" }},
			&Stringer{StringFunc: func() string { return "c" }},
		)

		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})
}
