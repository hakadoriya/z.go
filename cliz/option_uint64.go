package cliz

import (
	"slices"
	"strings"

	"github.com/hakadoriya/z.go/cliz/clicorez"
	"github.com/hakadoriya/z.go/errorz"
)

type (
	// Uint64Option is the option for int value.
	Uint64Option struct {
		// Name is the name of the option.
		Name string
		// Aliases is the alias names of the option.
		Aliases []string
		// Env is the environment variable name of the option.
		Env string
		// Default is the default value of the option.
		Default uint64
		// Required is the required flag of the option.
		Required bool
		// Description is the description of the option.
		Description string

		// value is the value of the option.
		value *uint64
	}
)

var _ Option = (*Uint64Option)(nil)

func (o *Uint64Option) GetName() string         { return o.Name }
func (o *Uint64Option) GetAliases() []string    { return o.Aliases }
func (o *Uint64Option) GetEnv() string          { return o.Env }
func (o *Uint64Option) GetDefault() interface{} { return o.Default }
func (o *Uint64Option) IsRequired() bool        { return o.Required }
func (o *Uint64Option) IsZero() bool            { return o.value == nil || *o.value == 0 }
func (o *Uint64Option) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "uint64 value of " + o.Name
}

func (c *Command) GetOptionUint64(name string) (uint64, error) {
	v, err := c.getOptionUint64(name)
	if err != nil {
		return 0, errorz.Errorf("cmd = %s: %w", strings.Join(c.allExecutedCommandNames, " "), err)
	}

	return v, nil
}

//nolint:cyclop
func (c *Command) getOptionUint64(name string) (uint64, error) {
	if len(c.allExecutedCommandNames) == 0 {
		return 0, errorz.Errorf("%s: %w", c.Name, clicorez.ErrNotCalled)
	}

	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range c.SubCommands {
		subcmd := c.SubCommands[len(c.SubCommands)-1-i]
		v, err := subcmd.getOptionUint64(name)
		if err == nil {
			return v, nil
		}
	}

	for _, opt := range c.Options {
		if o, ok := opt.(*Uint64Option); ok {
			if (o.Name != "" && o.Name == name) || (len(o.Aliases) > 0 && slices.Contains(o.Aliases, name)) || (o.Env != "" && o.Env == name) {
				if o.value != nil {
					return *o.value, nil
				}
			}
		}
	}

	return 0, errorz.Errorf("option = %s: %w", name, clicorez.ErrUnknownOption)
}
