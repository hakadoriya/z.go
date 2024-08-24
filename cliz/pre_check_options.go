package cliz

import (
	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/mapz"
)

func (c *Command) preCheckOptions() error {
	// NOTE: duplicate check
	if err := c.preCheckDuplicateOptions(make(map[string]duplicateMetadata)); err != nil {
		return errorz.Errorf("%s: %w", c.Name, err)
	}

	return nil
}

type duplicateMetadata struct {
	cmdName        string
	optionFullName string
}

//nolint:cyclop
func (c *Command) preCheckDuplicateOptions(duplicateChecker map[string]duplicateMetadata) error {
	for _, opt := range c.Options {
		if _, ok := opt.(*HelpOption); ok {
			// HelpOption is not a target of duplicate check
			continue
		}
		if name := opt.GetName(); name != "" {
			if key, alreadyExists := duplicateChecker[name]; alreadyExists {
				return errorz.Errorf("option: %s(%s) and %s(%s%s): %w", key.cmdName, key.optionFullName, c.Name, longOptionPrefix, name, ErrDuplicateOption)
			}
			duplicateChecker[name] = duplicateMetadata{
				cmdName:        c.Name,
				optionFullName: longOptionPrefix + name,
			}
		}

		for _, alias := range opt.GetAliases() {
			if alias != "" {
				if key, alreadyExists := duplicateChecker[alias]; alreadyExists {
					return errorz.Errorf("option: %s(%s) and %s(%s%s): %w", key.cmdName, key.optionFullName, c.Name, shortOptionPrefix, alias, ErrDuplicateOption)
				}
				duplicateChecker[alias] = duplicateMetadata{
					cmdName:        c.Name,
					optionFullName: shortOptionPrefix + alias,
				}
			}
		}

		if env := opt.GetEnv(); env != "" {
			const envPrefix = "env:"
			if key, alreadyExists := duplicateChecker[env]; alreadyExists {
				return errorz.Errorf("option: %s(%s) and %s(%s%s): %w", key.cmdName, key.optionFullName, c.Name, envPrefix, env, ErrDuplicateOption)
			}
			duplicateChecker[env] = duplicateMetadata{
				cmdName:        c.Name,
				optionFullName: envPrefix + env,
			}
		}
	}

	for _, subcmd := range c.SubCommands {
		// IMPORTANT: `copiedDuplicateChecker` is copied to avoid affecting the sibling subcommands.
		//            Ancestor command options and descendant command options must not be duplicated,
		//            but the sibling subcommands' options can be duplicated,
		//            because the sibling subcommands are not called at the same time,
		//            so they can coexist because sibling subcommands' options do not mix.
		copiedDuplicateChecker := mapz.Copy(duplicateChecker)
		if err := subcmd.preCheckDuplicateOptions(copiedDuplicateChecker); err != nil {
			return errorz.Errorf("%s: %w", subcmd.Name, err)
		}
	}

	return nil
}
