package cliz

import (
	"context"
	"errors"
	"testing"

	"github.com/hakadoriya/z.go/cliz/clicorez"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("error,ErrNilContext", func(t *testing.T) {
		t.Parallel()

		if _, err := FromContext(nil); !errors.Is(err, clicorez.ErrNilContext) {
			t.Errorf("❌: err != ErrNilContext: %+v", err)
		}
	})

	t.Run("error,ErrCommandNotSetInContext", func(t *testing.T) {
		t.Parallel()

		if _, err := FromContext(context.Background()); !errors.Is(err, clicorez.ErrNotFoundInContext) {
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

	t.Run("error,panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("❌: panic did not occur")

				if err, ok := r.(error); ok {
					t.Errorf("❌: err: %+v", err)

					if !errors.Is(err, clicorez.ErrNilContext) {
						t.Errorf("❌: err != ErrNilContext: %+v", err)
					}
				}
			}
		}()

		MustFromContext(nil)
	})
}
