package cliz

import (
	"context"
	"strings"
	"testing"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cmd, ok := FromContext(WithContext(context.Background(), &Command{}))
		if !ok {
			t.Fatalf("❌: !ok")
		}
		if cmd == nil {
			t.Errorf("❌: cmd == nil")
		}
	})
}

func TestMustFromContext(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := WithContext(context.Background(), &Command{})
		if c := MustFromContext(ctx); c == nil {
			t.Fatal("❌: c == nil")
		}
	})

	t.Run("error,panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("❌: panic did not occur")
			}

			err, ok := r.(error)
			if !ok {
				t.Fatalf("❌: err: %+v", err)
			}

			const expected = "runtime error: invalid memory address or nil pointer dereference"
			if err == nil || !strings.Contains(err.Error(), expected) {
				t.Fatalf("❌: err == nil || !strings.Contains(err.Error(), %q): %+v", expected, err)
			}
		}()

		MustFromContext(nil)
	})
}
