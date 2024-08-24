package contextz

import (
	"context"
	"errors"
	"testing"

	"github.com/hakadoriya/z.go/contextz/ctxerrz"
)

func TestWithValue(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		const expected = 1
		ctx := WithValue(context.Background(), expected)
		actual, err := Value[int](ctx)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}

		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})
}

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("error,ErrNilContext", func(t *testing.T) {
		t.Parallel()

		ctx := (context.Context)(nil)
		_, err := Value[int](ctx)
		if !errors.Is(err, ctxerrz.ErrNilContext) {
			t.Errorf("❌: !errors.Is(err, ctxerrz.ErrNilContext): %+v", err)
		}
	})

	t.Run("error,ErrNotFoundInContext", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, err := Value[int](ctx)
		if !errors.Is(err, ctxerrz.ErrNotFoundInContext) {
			t.Errorf("❌: !errors.Is(err, ctxerrz.ErrNotFoundInContext): %+v", err)
		}
	})
}

func TestMustValue(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		const expected = 1
		ctx := WithValue(context.Background(), expected)
		actual := MustValue[int](ctx)
		if expected != actual {
			t.Errorf("❌: expected(%v) != actual(%v)", expected, actual)
		}
	})

	t.Run("error,panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("❌: panic did not occur")

				if err, ok := r.(error); ok {
					t.Errorf("❌: err: %+v", err)

					if !errors.Is(err, ctxerrz.ErrNilContext) {
						t.Errorf("❌: !errors.Is(err, ctxerrz.ErrNilContext): %+v", err)
					}
				}
			}
		}()

		MustValue[int](nil)
	})
}
