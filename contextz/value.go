package contextz

import (
	"context"
	"fmt"

	"github.com/hakadoriya/z.go/contextz/ctxerrz"
)

type contextKey[T interface{}] struct{}

func WithValue[T interface{}](parent context.Context, val T) context.Context {
	return context.WithValue(parent, contextKey[T]{}, val)
}

func Value[T interface{}](ctx context.Context) (value T, err error) {
	if ctx == nil {
		return value, fmt.Errorf("%T: %w", contextKey[T]{}, ctxerrz.ErrNilContext)
	}

	val, ok := ctx.Value(contextKey[T]{}).(T)
	if !ok {
		return value, fmt.Errorf("%T: %w", contextKey[T]{}, ctxerrz.ErrNotFoundInContext)
	}

	return val, nil
}

func MustValue[T interface{}](ctx context.Context) T {
	val, err := Value[T](ctx)
	if err != nil {
		panic(fmt.Errorf("Value: %w", err))
	}

	return val
}
