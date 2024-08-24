package cliz

import (
	"github.com/hakadoriya/z.go/errorz"
)

func (c *Command) postCheckOptions() error {
	// NOTE: required options
	if err := c.postCheckOptionRequired(); err != nil {
		return errorz.Errorf("%s: %w", c.Name, err)
	}

	return nil
}

func (c *Command) postCheckOptionRequired() error {
	if len(c.allExecutedCommandNames) > 0 {
		for _, opt := range c.Options {
			name := opt.GetName()
			if opt.IsRequired() && opt.IsZero() {
				return errorz.Errorf("%s: option: %s%s: %w", c.Name, longOptionPrefix, name, ErrOptionRequired)
			}
		}
	}

	for _, subcmd := range c.SubCommands {
		if err := subcmd.postCheckOptionRequired(); err != nil {
			return errorz.Errorf("%s: %w", subcmd.Name, err)
		}
	}

	return nil
}
