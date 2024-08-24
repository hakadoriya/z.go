package cliz

import (
	"context"

	"github.com/hakadoriya/z.go/contextz"
)

func (c *Command) Context() context.Context {
	return c.ctx
}

func (c *Command) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func WithContext(ctx context.Context, cmd *Command) context.Context {
	return contextz.WithValue(ctx, cmd)
}

func FromContext(ctx context.Context) (*Command, error) {
	return contextz.Value[*Command](ctx)
}

func MustFromContext(ctx context.Context) *Command {
	return contextz.MustValue[*Command](ctx)
}
