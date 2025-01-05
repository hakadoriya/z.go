package contextz

import (
	"context"
	"fmt"
)

func WithValue[T interface{}](parent context.Context, val T) context.Context {
	return context.WithValue(parent, (*T)(nil), val)
}

func Value[T interface{}](ctx context.Context) (value T, err error) {
	if ctx == nil {
		return value, fmt.Errorf("%T: %w", (*T)(nil), ErrNilContext)
	}

	val, ok := ctx.Value((*T)(nil)).(T)
	if !ok {
		return value, fmt.Errorf("%T: %w", (*T)(nil), ErrNotFoundInContext)
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
