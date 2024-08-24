package contextz

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

type testContext struct {
	DeadlineFunc func() (deadline time.Time, ok bool)
	DoneFunc     func() <-chan struct{}
	ErrFunc      func() error
	ValueFunc    func(key interface{}) interface{}
}

var _ context.Context = (*testContext)(nil)

func (c *testContext) Deadline() (deadline time.Time, ok bool) { return c.DeadlineFunc() }
func (c *testContext) Done() <-chan struct{}                   { return c.DoneFunc() }
func (c *testContext) Err() error                              { return c.ErrFunc() }
func (c *testContext) Value(key interface{}) interface{}       { return c.ValueFunc(key) }

func TestCheckContextDone(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		err := CheckContext(context.Background())
		if err != nil {
			t.Fatalf("❌: err != nil: %v", err)
		}
	})

	t.Run("error,context.Canceled", func(t *testing.T) {
		t.Parallel()

		contextCanceled, cancelCause := context.WithCancelCause(context.Background())
		cancelCause(nil)
		err := CheckContext(contextCanceled)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("❌: !errors.Is(err, context.Canceled): %v", err)
		}
	})

	t.Run("error,io.ErrUnexpectedEOF", func(t *testing.T) {
		t.Parallel()

		contextCanceled, cancelCause := context.WithCancelCause(context.Background())
		cancelCause(io.ErrUnexpectedEOF)
		err := CheckContext(contextCanceled)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("❌: !errors.Is(err, io.ErrUnexpectedEOF): %v", err)
		}
	})

	t.Run("error,nil", func(t *testing.T) {
		t.Parallel()

		errCalledCount := 0
		err := CheckContext(&testContext{
			DoneFunc: func() <-chan struct{} {
				closed := make(chan struct{})
				close(closed)
				return closed
			},
			ErrFunc: func() error {
				errCalledCount++
				if errCalledCount == 1 {
					return nil
				}
				return context.Canceled
			},
			ValueFunc: func(key interface{}) interface{} {
				return nil
			},
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("❌: !errors.Is(err, contextCanceled): %v", err)
		}
	})
}
