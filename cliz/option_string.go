package cliz

import (
	"slices"
	"strings"

	"github.com/hakadoriya/z.go/cliz/clierrz"
	"github.com/hakadoriya/z.go/errorz"
)

type (
	// StringOption is the option for string value.
	StringOption struct {
		// Name is the name of the option.
		Name string
		// Aliases is the alias names of the option.
		Aliases []string
		// Environment is the environment variable name of the option.
		Environment string
		// Default is the default value of the option.
		Default string
		// Required is the required flag of the option.
		Required bool
		// Description is the description of the option.
		Description string

		// value is the value of the option.
		value *string
	}
)

var _ Option = (*StringOption)(nil)

func (o *StringOption) GetName() string         { return o.Name }
func (o *StringOption) GetAliases() []string    { return o.Aliases }
func (o *StringOption) GetEnvironment() string  { return o.Environment }
func (o *StringOption) GetDefault() interface{} { return o.Default }
func (o *StringOption) IsRequired() bool        { return o.Required }
func (o *StringOption) IsZero() bool            { return o.value == nil || *o.value == "" }
func (o *StringOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "string value of " + o.Name
}

func (c *Command) GetOptionString(name string) (string, error) {
	v, err := c.getOptionString(name)
	if err != nil {
		return "", errorz.Errorf("cmd = %s: %w", strings.Join(c.allExecutedCommandNames, " "), err)
	}

	return v, nil
}

//nolint:cyclop
func (c *Command) getOptionString(name string) (string, error) {
	if len(c.allExecutedCommandNames) == 0 {
		return "", errorz.Errorf("%s: %w", c.Name, clierrz.ErrNotCalled)
	}

	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range c.SubCommands {
		subcmd := c.SubCommands[len(c.SubCommands)-1-i]
		v, err := subcmd.getOptionString(name)
		if err == nil {
			return v, nil
		}
	}

	for _, opt := range c.Options {
		if o, ok := opt.(*StringOption); ok {
			if (o.Name != "" && o.Name == name) || (len(o.Aliases) > 0 && slices.Contains(o.Aliases, name)) || (o.Environment != "" && o.Environment == name) {
				if o.value != nil {
					return *o.value, nil
				}
			}
		}
	}

	return "", errorz.Errorf("option = %s: %w", name, clierrz.ErrUnknownOption)
}
