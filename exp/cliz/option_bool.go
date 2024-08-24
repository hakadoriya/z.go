package cliz

import (
	"slices"

	"github.com/hakadoriya/z.go/errorz"
)

type (
	// BoolOption is the option for bool value.
	BoolOption struct {
		// Name is the name of the option.
		Name string
		// Aliases is the alias names of the option.
		Aliases []string
		// Environment is the environment variable name of the option.
		Environment string
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
func (o *BoolOption) GetEnvironment() string  { return o.Environment }
func (o *BoolOption) GetDefault() interface{} { return o.Default }
func (o *BoolOption) IsRequired() bool        { return o.Required }
func (o *BoolOption) IsZero() bool            { return o.value == nil || *o.value == false }
func (o *BoolOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "bool value of " + o.Name
}

func (cmd *Command) GetOptionBool(name string) (bool, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetOptionBool(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getOptionBool(name)
	if err == nil {
		return v, nil
	}

	return false, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

func (cmd *Command) getOptionBool(name string) (bool, error) {
	if len(cmd.calledCommands) == 0 {
		return false, errorz.Errorf("%s: %w", cmd.Name, ErrNotCalled)
	}

	for _, opt := range cmd.Options {
		if o, ok := opt.(*BoolOption); ok {
			if (o.Name != "" && o.Name == name) || (len(o.Aliases) > 0 && slices.Contains(o.Aliases, name)) || (o.Environment != "" && o.Environment == name) {
				if o.value != nil {
					return *o.value, nil
				}
			}
		}
	}

	return false, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}
