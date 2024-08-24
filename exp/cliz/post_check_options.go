package cliz

import (
	"github.com/hakadoriya/z.go/errorz"
)

func (cmd *Command) postCheckOptions() error {
	// NOTE: required options
	if err := cmd.postCheckOptionRequired(); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	if err := cmd.postCheckDuplicateOptionInCalledCommands(make(map[string]bool), make(map[string]bool)); err != nil {
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

func (cmd *Command) postCheckDuplicateOptionInCalledCommands(names map[string]bool, aliases map[string]bool) error {
	if cmd.called {
		for _, opt := range cmd.Options {
			if _, ok := opt.(*HelpOption); ok {
				continue
			}
			if name := opt.GetName(); name != "" {
				if names[name] {
					return errorz.Errorf("option: %s%s: %w", longOptionPrefix, name, ErrDuplicateOption)
				}
				names[name] = true
			}

			for _, alias := range opt.GetAliases() {
				if alias != "" {
					if aliases[alias] {
						return errorz.Errorf("short option: %s%s: %w", shortOptionPrefix, alias, ErrDuplicateOption)
					}
					aliases[alias] = true
				}
			}
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.postCheckDuplicateOptionInCalledCommands(names, aliases); err != nil {
			return errorz.Errorf("subcommand: %s: %w", subcmd.Name, err)
		}
	}

	return nil
}
