package cliz

import (
	"slices"
	"strings"

	"github.com/hakadoriya/z.go/cliz/clicorez"
	"github.com/hakadoriya/z.go/errorz"
)

type (
	// BoolOption is the option for bool value.
	BoolOption struct {
		// Name is the name of the option.
		Name string
		// Aliases is the alias names of the option.
		Aliases []string
		// Env is the environment variable name of the option.
		Env string
		// Default is the default value of the option.
		Default bool
		// Required is the required flag of the option.
		Required bool
		// Description is the description of the option.
		Description string

		// value is the value of the option.
		value *bool
	}
)

var _ Option = (*BoolOption)(nil)

func (o *BoolOption) GetName() string         { return o.Name }
func (o *BoolOption) GetAliases() []string    { return o.Aliases }
func (o *BoolOption) GetEnv() string          { return o.Env }
func (o *BoolOption) GetDefault() interface{} { return o.Default }
func (o *BoolOption) IsRequired() bool        { return o.Required }
func (o *BoolOption) IsZero() bool            { return o.value == nil || !*o.value }
func (o *BoolOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "bool value of " + o.Name
}

func (c *Command) GetOptionBool(name string) (bool, error) {
	v, err := c.getOptionBool(name)
	if err != nil {
		return false, errorz.Errorf("cmd = %s: %w", strings.Join(c.allExecutedCommandNames, " "), err)
	}

	return v, nil
}

//nolint:cyclop
func (c *Command) getOptionBool(name string) (bool, error) {
	if len(c.allExecutedCommandNames) == 0 {
		return false, errorz.Errorf("%s: %w", c.Name, clicorez.ErrNotCalled)
	}

	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range c.SubCommands {
		subcmd := c.SubCommands[len(c.SubCommands)-1-i]
		v, err := subcmd.getOptionBool(name)
		if err == nil {
			return v, nil
		}
	}

	for _, opt := range c.Options {
		if o, ok := opt.(*BoolOption); ok {
			if (o.Name != "" && o.Name == name) || (len(o.Aliases) > 0 && slices.Contains(o.Aliases, name)) || (o.Env != "" && o.Env == name) {
				if o.value != nil {
					return *o.value, nil
				}
			}
		}
	}

	return false, errorz.Errorf("option = %s: %w", name, clicorez.ErrUnknownOption)
}
