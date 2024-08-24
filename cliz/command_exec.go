package cliz

import (
	"context"
	"strings"

	"github.com/hakadoriya/z.go/cliz/clierrz"
	"github.com/hakadoriya/z.go/errorz"
)

func (c *Command) Exec(ctx context.Context, osArgs []string) error {
	remainingArgs, err := c.parse(ctx, osArgs)
	if err != nil {
		return errorz.Errorf("cmd = %s: c.parse: %w", strings.Join(c.allExecutedCommandNames, " "), err)
	}

	executed := c.GetExecutedCommand()

	if executed.PreHookExecFunc != nil {
		if err := executed.PreHookExecFunc(ctx, c, remainingArgs); err != nil {
			return errorz.Errorf("cmd = %s: called.PreHookExecFunc: %w", strings.Join(c.allExecutedCommandNames, " "), err)
		}
	}

	if executed.ExecFunc == nil {
		executed.showUsage()
		return clierrz.ErrHelp
	}

	if err := executed.ExecFunc(ctx, c, remainingArgs); err != nil {
		return errorz.Errorf("%s: ExecFunc: %w", executed.Name, err)
	}

	if executed.PostHookExecFunc != nil {
		if err := executed.PostHookExecFunc(ctx, c, remainingArgs); err != nil {
			return errorz.Errorf("%s: PostHookFunc: %w", executed.Name, err)
		}
	}

	return nil
}
