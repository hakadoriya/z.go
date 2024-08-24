package cliz

import "github.com/hakadoriya/z.go/errorz"

func (cmd *Command) preCheckOptions() error {
	// NOTE: duplicate check
	if err := cmd.preCheckDuplicateOptions(); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	return nil
}

//nolint:cyclop
func (cmd *Command) preCheckDuplicateOptions() error {
	envs := make(map[string]bool)
	names := make(map[string]bool)
	aliases := make(map[string]bool)

	for _, opt := range cmd.Options {
		if name := opt.GetName(); name != "" {
			if names[name] {
				err := ErrDuplicateOption
				return errorz.Errorf("option: %s%s: %w", longOptionPrefix, name, err)
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

		if env := opt.GetEnvironment(); env != "" {
			if envs[env] {
				return errorz.Errorf("environment: %s: %w", env, ErrDuplicateOption)
			}
			envs[env] = true
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.preCheckDuplicateOptions(); err != nil {
			return errorz.Errorf("subcommand: %s: %w", subcmd.Name, err)
		}
	}

	return nil
}
