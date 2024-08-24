package cliz

import (
	"context"

	"github.com/hakadoriya/z.go/errorz"
)

func WithContext(ctx context.Context, cmd *Command) context.Context {
	return context.WithValue(ctx, (*Command)(nil), cmd)
}

func FromContext(ctx context.Context) (*Command, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	c, ok := ctx.Value((*Command)(nil)).(*Command)
	if !ok {
		return nil, errorz.Errorf("%T: %w", (*Command)(nil), ErrNotSetInContext)
	}

	return c, nil
}

func MustFromContext(ctx context.Context) *Command {
	c, err := FromContext(ctx)
	if err != nil {
		err = errorz.Errorf("FromContext: %w", err)
		panic(err)
	}

	return c
}
