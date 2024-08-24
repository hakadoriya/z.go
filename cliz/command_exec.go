package cliz

import (
	"context"
	"strings"

	"github.com/hakadoriya/z.go/cliz/clicorez"
	"github.com/hakadoriya/z.go/errorz"
)

func (c *Command) Exec(ctx context.Context, osArgs []string) (err error) {
	remainingArgs, err := c.parse(ctx, osArgs)
	if err != nil {
		return errorz.Errorf("cmd = %s: parse: %w", strings.Join(c.allExecutedCommandNames, " "), err)
	}

	executed := c.GetExecutedCommand()

	defer func() {
		if IsHelp(err) {
			executed.showUsage()
		}
	}()

	if executed.PreHookExecFunc != nil {
		if err := executed.PreHookExecFunc(c, remainingArgs); err != nil {
			return errorz.Errorf("cmd = %s: PreHookExec: %w", strings.Join(c.allExecutedCommandNames, " "), err)
		}
	}

	if executed.ExecFunc == nil {
		return clicorez.ErrHelp
	}

	if err := executed.ExecFunc(c, remainingArgs); err != nil {
		return errorz.Errorf("cmd = %s: Exec: %w", strings.Join(c.allExecutedCommandNames, " "), err)
	}

	if executed.PostHookExecFunc != nil {
		if err := executed.PostHookExecFunc(c, remainingArgs); err != nil {
			return errorz.Errorf("cmd = %s: PostHookExec: %w", strings.Join(c.allExecutedCommandNames, " "), err)
		}
	}

	return nil
}
