package cliz

import (
	"github.com/hakadoriya/z.go/cliz/clierrz"
	"github.com/hakadoriya/z.go/errorz"
)

func (c *Command) preCheckSubCommands() error {
	if err := c.preCheckDuplicateSubCommands(); err != nil {
		return errorz.Errorf("%s: %w", c.Name, err)
	}

	return nil
}

func (c *Command) preCheckDuplicateSubCommands() error {
	names := make(map[string]bool)

	for _, subcmd := range c.SubCommands {
		name := subcmd.Name

		if name != "" && names[name] {
			return errorz.Errorf("subcommand: %s: %w", name, clierrz.ErrDuplicateSubCommand)
		}
		names[name] = true
	}

	for _, subcmd := range c.SubCommands {
		if err := subcmd.preCheckDuplicateSubCommands(); err != nil {
			return errorz.Errorf("subcommand: %s: %w", subcmd.Name, err)
		}
	}

	return nil
}
