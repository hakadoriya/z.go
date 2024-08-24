package cliz

import "github.com/hakadoriya/z.go/errorz"

// Default is the helper function to create a default value.
func Default[T interface{}](v T) *T { return ptr[T](v) }

func ptr[T interface{}](v T) *T { return &v }

//nolint:cyclop
func (cmd *Command) loadDefaults() error {
	for _, opt := range cmd.Options {
		switch o := opt.(type) {
		case *StringOption:
			o.value = &o.Default
		case *BoolOption:
			o.value = &o.Default
		case *Int64Option:
			o.value = &o.Default
		case *Float64Option:
			o.value = &o.Default
		default:
			return errorz.Errorf("%s: %w", o.GetName(), ErrInvalidOptionType)
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.loadDefaults(); err != nil {
			return err
		}
	}

	return nil
}
