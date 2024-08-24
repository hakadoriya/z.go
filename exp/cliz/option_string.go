package cliz

import (
	"slices"

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

func (cmd *Command) GetOptionString(name string) (string, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetOptionString(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getOptionString(name)
	if err == nil {
		return v, nil
	}

	return "", errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

func (cmd *Command) getOptionString(name string) (string, error) {
	if len(cmd.calledCommands) == 0 {
		return "", errorz.Errorf("%s: %w", cmd.Name, ErrNotCalled)
	}

	for _, opt := range cmd.Options {
		if o, ok := opt.(*StringOption); ok {
			if (o.Name != "" && o.Name == name) || (len(o.Aliases) > 0 && slices.Contains(o.Aliases, name)) || (o.Environment != "" && o.Environment == name) {
				if o.value != nil {
					return *o.value, nil
				}
			}
		}
	}

	return "", errorz.Errorf("%s: %w", name, ErrUnknownOption)
}
