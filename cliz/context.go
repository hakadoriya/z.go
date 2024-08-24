package cliz

import (
	"context"
)

func (c *Command) Context() context.Context {
	return c.ctx
}

func (c *Command) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func WithContext(ctx context.Context, cmd *Command) context.Context {
	return context.WithValue(ctx, (*Command)(nil), cmd)
}

func FromContext(ctx context.Context) (cmd *Command, ok bool) {
	cmd, ok = ctx.Value((*Command)(nil)).(*Command)
	return cmd, ok
}

func MustFromContext(ctx context.Context) *Command {
	// panic intentionally if the value is not a *Command.
	//
	//nolint:forcetypeassert
	return ctx.Value((*Command)(nil)).(*Command)
}
