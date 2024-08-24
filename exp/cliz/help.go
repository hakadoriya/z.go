package cliz

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

// HelpOptionName is the option name for help.
const HelpOptionName = "help"

// IsHelp returns whether the error is ErrHelp.
func IsHelp(err error) bool {
	return errors.Is(err, ErrHelp)
}

func (cmd *Command) initAppendHelpOption() {
	// If help option is already set, do nothing.
	if _, ok := cmd.getHelpOption(); !ok {
		cmd.Options = append(cmd.Options, &BoolOption{
			Name:        HelpOptionName,
			Description: "show usage",
			Default:     false,
			Required:    false,
		})
	}

	// Recursively initialize help option for subcommands.
	for _, subcmd := range cmd.SubCommands {
		subcmd.initAppendHelpOption()
	}
}

func (cmd *Command) getHelpOption() (helpOption *BoolOption, ok bool) {
	// Find help option in the command options.
	for _, opt := range cmd.Options {
		if o, ok := opt.(*BoolOption); ok {
			if o.Name == HelpOptionName {
				return o, true
			}
		}
	}

	return nil, false
}

func (cmd *Command) checkHelp(ctx context.Context) error {
	Logger.Debug("checkHelp: " + cmd.Name)

	// If help option is set, show usage and return ErrHelp.
	v, err := cmd.getOptionBool(HelpOptionName)
	if err == nil && v {
		cmd.ShowUsage(stderrFromContext(ctx))
		return ErrHelp
	}

	// Recursively check help option for subcommands.
	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.checkHelp(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *Command) ShowUsage(stderr io.Writer) {
	if cmd.UsageFunc != nil {
		cmd.UsageFunc(stderr)
		return
	}
	showUsage(stderr, cmd)
}

//nolint:cyclop,funlen,gocognit
func showUsage(stderr io.Writer, cmd *Command) {
	const indent = "    "

	// Usage
	usage := "Usage:" + "\n"
	if cmd.Usage != "" {
		usage += indent + cmd.Usage + "\n"
	} else {
		usage += indent + strings.Join(cmd.GetCalledCommands(), " ")
		if len(cmd.Options) > 0 {
			usage += " [options]"
		}
		if len(cmd.SubCommands) > 0 {
			usage += " <subcommand>"
		}
		usage += "\n"
	}
	usage += "\n"

	// Description
	usage += "Description:" + "\n"
	usage += indent + cmd.Description + "\n"

	// SubCommands
	if len(cmd.SubCommands) > 0 {
		usage += "\n"
		usage += "Sub commands:\n"
		groups := make([]string, 0)
		subCommandsByGroup := make(map[string][]*Command)
		for _, subcmd := range cmd.SubCommands {
			groups = append(groups, subcmd.Group)
			subCommandsByGroup[subcmd.Group] = append(subCommandsByGroup[subcmd.Group], subcmd)
		}
		slices.Sort(groups)
		for _, group := range slices.Compact(groups) {
			// If the group is set, the group name is displayed.
			if group != "" {
				usage += indent + group + ":\n"
			}
			var commandsMaxWidthInGroup int
			for _, subcmd := range subCommandsByGroup[group] {
				if len(subcmd.Name) > commandsMaxWidthInGroup {
					commandsMaxWidthInGroup = len(subcmd.Name)
				}
			}
			for _, subcmd := range subCommandsByGroup[group] {
				// If the group is set, add an indent for group name.
				if subcmd.Group != "" {
					usage += indent
				}
				usage += indent + fmt.Sprintf("%-"+strconv.Itoa(commandsMaxWidthInGroup)+"s"+indent+"%s", subcmd.Name, subcmd.Description) + "\n"
			}
		}
	}

	// Options
	if len(cmd.Options) > 0 { //nolint:nestif
		usage += "\n"
		usage += "Options:\n"
		for _, opt := range cmd.Options {
			name := opt.GetName()
			env := opt.GetEnvironment()
			usage += indent
			if name != "" {
				usage += fmt.Sprintf("%s%s", longOptionPrefix, name)
			}
			if aliases := opt.GetAliases(); len(aliases) > 0 {
				for _, alias := range aliases {
					usage += fmt.Sprintf(", %s%s", shortOptionPrefix, alias)
				}
			}

			usage += " ("

			if env != "" {
				usage += fmt.Sprintf("env: %s, ", env)
			}

			usage += fmt.Sprintf("default: %v", opt.GetDefault())

			usage += ")"

			usage += "\n"
			usage += indent + indent + opt.GetDescription() + "\n"
		}
	}

	// Output
	_, _ = fmt.Fprint(stderr, usage)
}
