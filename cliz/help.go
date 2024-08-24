package cliz

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hakadoriya/z.go/slicez"
)

// HelpOptionName is the option name for help.
const HelpOptionName = "help"

// IsHelp returns whether the error is ErrHelp.
func IsHelp(err error) bool {
	return errors.Is(err, ErrHelp)
}

func (c *Command) initAppendHelpOption() {
	// If help option is already set, do nothing.
	if _, ok := c.getHelpOption(); !ok {
		//nolint:exhaustruct
		c.Options = append(c.Options, &HelpOption{
			Name: HelpOptionName,
		})
	}

	// Recursively initialize help option for subcommands.
	for _, subcmd := range c.SubCommands {
		subcmd.initAppendHelpOption()
	}
}

func (c *Command) getHelpOption() (helpOption *HelpOption, ok bool) {
	// Find help option in the command options.
	for _, opt := range c.Options {
		if o, ok := opt.(*HelpOption); ok {
			if o.Name == HelpOptionName {
				return o, true
			}
		}
	}

	return nil, false
}

func (c *Command) checkHelp() error {
	// If help option is set, show help message and return ErrHelp.
	helpRequested, err := c.getOptionHelp(HelpOptionName)
	if err == nil && helpRequested {
		Logger.Debug("checkHelp: " + strings.Join(c.allExecutedCommandNames, " "))
		c.showUsage()
		return ErrHelp
	}

	// Recursively check help option for subcommands.
	for _, subcmd := range c.SubCommands {
		if err := subcmd.checkHelp(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) showUsage() {
	if c.UsageFunc != nil {
		c.UsageFunc(c)
		return
	}

	c.DefaultUsage()
}

//nolint:cyclop,funlen,gocognit
func (c *Command) DefaultUsage() {
	Logger.Debug("DefaultUsage: " + c.Name)

	const indent = "    "

	// Usage
	usage := "Usage:" + "\n"
	if c.Usage != "" {
		usage += indent + c.Usage + "\n"
	} else {
		usage += indent + strings.Join(c.allExecutedCommandNames, " ")
		if len(c.Options) > 0 {
			usage += " [options]"
		}
		if len(c.SubCommands) > 0 && !c.hasOnlyHiddenSubCommands() {
			usage += " <subcommand>"
		}
		usage += "\n"
	}

	// Description
	if c.Description != "" {
		usage += "\n"
		usage += "Description:" + "\n"
		usage += indent + c.Description + "\n"
	}

	// SubCommands
	if len(c.SubCommands) > 0 && !c.hasOnlyHiddenSubCommands() {
		usage += "\n"
		usage += "Sub Commands:\n"
		groups := c.getGroups()
		subCommandsByGroup := c.getSubCommandsByGroup()
		for _, group := range groups {
			// If the group is set, the group name is displayed.
			if group != "" {
				usage += indent + group + ":\n"
			}
			commandsMaxWidthInGroup := c.getSubCommandListMaxWidthInGroup(group)
			for _, subcmd := range subCommandsByGroup[group] {
				// If the group is set, add an indent for group name.
				if subcmd.Group != "" {
					usage += indent
				}
				usage += indent + fmt.Sprintf("%-"+strconv.Itoa(commandsMaxWidthInGroup)+"s"+indent+"%s", subcmd.getNameAndAliasesString(), subcmd.Description) + "\n"
			}
		}
	}

	// Options
	//
	//nolint:nestif
	if len(c.Options) > 0 {
		usage += "\n"
		usage += "Options:\n"
		for _, opt := range c.Options {
			if opt.IsHidden() {
				continue
			}

			name := opt.GetName()
			env := opt.GetEnv()
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

			if opt.IsRequired() {
				usage += "required, "
			}

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
	_, _ = io.WriteString(c.Stderr(), usage)
}

func (c *Command) getNameAndAliasesString() string {
	names := make([]string, 0)
	names = append(names, c.Name)
	names = append(names, c.Aliases...)
	return strings.Join(slicez.CompactStable(names), ", ")
}

func (c *Command) getGroups() (groups []string) {
	groups = make([]string, 0)
	for _, subcmd := range c.SubCommands {
		if subcmd.Hidden {
			continue
		}
		groups = append(groups, subcmd.Group)
	}
	return slicez.CompactStable(groups)
}

func (c *Command) getSubCommandsByGroup() map[string][]*Command {
	subCommandsByGroup := make(map[string][]*Command)
	for _, subcmd := range c.SubCommands {
		if subcmd.Hidden {
			continue
		}
		subCommandsByGroup[subcmd.Group] = append(subCommandsByGroup[subcmd.Group], subcmd)
	}
	return subCommandsByGroup
}

func (c *Command) getSubCommandListMaxWidthInGroup(group string) (maxWidth int) {
	for _, subcmd := range c.SubCommands {
		if subcmd.Hidden || subcmd.Group != group {
			continue
		}
		nameAndAliasesString := subcmd.getNameAndAliasesString()
		if len(nameAndAliasesString) > maxWidth {
			maxWidth = len(nameAndAliasesString)
		}
	}
	return maxWidth
}
