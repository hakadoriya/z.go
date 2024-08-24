package cliz

import "github.com/hakadoriya/z.go/errorz"

func (cmd *Command) preCheckSubCommands() error {
	if err := cmd.preCheckDuplicateSubCommands(); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	return nil
}

func (cmd *Command) preCheckDuplicateSubCommands() error {
	names := make(map[string]bool)

	for _, cmd := range cmd.SubCommands {
		name := cmd.Name

		if name != "" && names[name] {
			return errorz.Errorf("subcommand: %s: %w", name, ErrDuplicateSubCommand)
		}
		names[name] = true
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.preCheckDuplicateSubCommands(); err != nil {
			return errorz.Errorf("subcommand: %s: %w", subcmd.Name, err)
		}
	}

	return nil
}
