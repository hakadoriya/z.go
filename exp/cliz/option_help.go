package cliz

import (
	"slices"

	"github.com/hakadoriya/z.go/errorz"
)

type (
	// HelpOption is the option for bool value.
	HelpOption struct {
		// Name is the name of the option.
		Name string
		// Aliases is the alias names of the option.
		Aliases []string
		// Description is the description of the option.
		Description string

		// value is the value of the option.
		value *bool
	}
)

var _ Option = (*HelpOption)(nil)

func (o *HelpOption) GetName() string         { return o.Name }
func (o *HelpOption) GetAliases() []string    { return o.Aliases }
func (o *HelpOption) GetEnvironment() string  { return "" }
func (o *HelpOption) GetDefault() interface{} { return false }
func (o *HelpOption) IsRequired() bool        { return false }
func (o *HelpOption) IsZero() bool            { return o.value == nil || *o.value == false }
func (o *HelpOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "show usage"
}

func (cmd *Command) getOptionHelp(name string) (bool, error) {
	if len(cmd.calledCommands) == 0 {
		return false, errorz.Errorf("%s: %w", cmd.Name, ErrNotCalled)
	}

	for _, opt := range cmd.Options {
		if o, ok := opt.(*HelpOption); ok {
			if (o.Name != "" && o.Name == name) || (len(o.Aliases) > 0 && slices.Contains(o.Aliases, name)) {
				if o.value != nil {
					return *o.value, nil
				}
			}
		}
	}

	return false, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}
