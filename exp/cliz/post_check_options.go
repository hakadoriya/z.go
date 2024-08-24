package cliz

import (
	"github.com/hakadoriya/z.go/errorz"
)

func (cmd *Command) postCheckOptions() error {
	// NOTE: required options
	if err := cmd.postCheckOptionRequired(); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	return nil
}

//nolint:cyclop
func (cmd *Command) postCheckOptionRequired() error {
	if cmd.called {
		for _, opt := range cmd.Options {
			name := opt.GetName()
			if opt.IsRequired() && opt.IsZero() {
				return errorz.Errorf("%s: option: %s%s: %w", cmd.Name, longOptionPrefix, name, ErrOptionRequired)
			}
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.postCheckOptionRequired(); err != nil {
			return errorz.Errorf("%s: %w", subcmd.Name, err)
		}
	}

	return nil
}
