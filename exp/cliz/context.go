package cliz

import (
	"context"
	"io"

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

type (
	contextKeyStdout struct{}
	contextKeyStderr struct{}
)

func stdoutWithContext(ctx context.Context, s io.Writer) context.Context {
	return context.WithValue(ctx, contextKeyStdout{}, s)
}

func stdoutFromContext(ctx context.Context) io.Writer {
	s, ok := ctx.Value(contextKeyStdout{}).(io.Writer)
	if !ok {
		panic(errorz.Errorf("%T: %w", contextKeyStdout{}, ErrNotSetInContext))

	}

	return s
}

func stderrWithContext(ctx context.Context, s io.Writer) context.Context {
	return context.WithValue(ctx, contextKeyStderr{}, s)
}

func stderrFromContext(ctx context.Context) io.Writer {
	s, ok := ctx.Value(contextKeyStderr{}).(io.Writer)
	if !ok {
		panic(errorz.Errorf("%T: %w", contextKeyStderr{}, ErrNotSetInContext))

	}

	return s
}
