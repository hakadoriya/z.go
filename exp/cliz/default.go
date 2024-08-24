package cliz

import "github.com/hakadoriya/z.go/errorz"

func ptr[T interface{}](v T) *T { return &v }

//nolint:cyclop
func (cmd *Command) loadDefaults() error {
	for _, opt := range cmd.Options {
		switch o := opt.(type) {
		case *StringOption:
			if !o.IsRequired() {
				o.value = &o.Default
			}
		case *BoolOption:
			if !o.IsRequired() {
				o.value = &o.Default
			}
		case *Int64Option:
			if !o.IsRequired() {
				o.value = &o.Default
			}
		case *Float64Option:
			if !o.IsRequired() {
				o.value = &o.Default
			}
		case *HelpOption:
			// do nothing
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
