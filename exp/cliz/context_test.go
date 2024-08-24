package cliz

import (
	"context"
	"errors"
	"testing"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("failure,ErrNilContext", func(t *testing.T) {
		t.Parallel()

		if _, err := FromContext(nil); !errors.Is(err, ErrNilContext) {
			t.Errorf("❌: err != ErrNilContext: %+v", err)
		}
	})

	t.Run("failure,ErrCommandNotSetInContext", func(t *testing.T) {
		t.Parallel()

		if _, err := FromContext(context.Background()); !errors.Is(err, ErrNotSetInContext) {
			t.Errorf("❌: err != ErrCommandNotSetInContext: %+v", err)
		}
	})
}

func TestMustFromContext(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := WithContext(context.Background(), &Command{})
		if c := MustFromContext(ctx); c == nil {
			t.Errorf("❌: c == nil")
		}
	})

	t.Run("failure,panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("❌: panic did not occur")

				if err, ok := r.(error); ok {
					t.Errorf("❌: err: %+v", err)

					if !errors.Is(err, ErrNilContext) {
						t.Errorf("❌: err != ErrNilContext: %+v", err)
					}
				}
			}
		}()

		MustFromContext(nil)
	})
}
