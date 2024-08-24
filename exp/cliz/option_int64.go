package cliz

import (
	"slices"

	"github.com/hakadoriya/z.go/errorz"
)

type (
	// Int64Option is the option for int value.
	Int64Option struct {
		// Name is the name of the option.
		Name string
		// Aliases is the alias names of the option.
		Aliases []string
		// Environment is the environment variable name of the option.
		Environment string
		// Default is the default value of the option.
		Default int64
		// Required is the required flag of the option.
		Required bool
		// Description is the description of the option.
		Description string

		// value is the value of the option.
		value *int64
	}
)

var _ Option = (*Int64Option)(nil)

func (o *Int64Option) GetName() string         { return o.Name }
func (o *Int64Option) GetAliases() []string    { return o.Aliases }
func (o *Int64Option) GetEnvironment() string  { return o.Environment }
func (o *Int64Option) GetDefault() interface{} { return o.Default }
func (o *Int64Option) IsRequired() bool        { return o.Required }
func (o *Int64Option) IsZero() bool            { return o.value == nil || *o.value == 0 }
func (o *Int64Option) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "int64 value of " + o.Name
}

func (cmd *Command) GetOptionInt64(name string) (int64, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetOptionInt64(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getOptionInt64(name)
	if err == nil {
		return v, nil
	}

	return 0, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

func (cmd *Command) getOptionInt64(name string) (int64, error) {
	if len(cmd.calledCommands) == 0 {
		return 0, errorz.Errorf("%s: %w", cmd.Name, ErrNotCalled)
	}

	for _, opt := range cmd.Options {
		if o, ok := opt.(*Int64Option); ok {
			if (o.Name != "" && o.Name == name) || (len(o.Aliases) > 0 && slices.Contains(o.Aliases, name)) || (o.Environment != "" && o.Environment == name) {
				if o.value != nil {
					return *o.value, nil
				}
			}
		}
	}

	return 0, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}
