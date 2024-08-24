package cliz

import (
	"os"
	"strconv"

	"github.com/hakadoriya/z.go/errorz"
)

//nolint:gocognit,cyclop
func (cmd *Command) loadEnvironments() error {
	for _, opt := range cmd.Options {
		if opt.GetEnvironment() == "" {
			// If v is an empty string, o.value remains.
			continue
		}

		switch o := opt.(type) {
		case *StringOption:
			if s := os.Getenv(o.Environment); s != "" {
				o.value = &s
			}
		case *BoolOption:
			if s := os.Getenv(o.Environment); s != "" {
				v, err := strconv.ParseBool(s)
				if err != nil {
					return errorz.Errorf("%s: %w", o.Environment, err)
				}
				o.value = &v
			}
		case *Int64Option:
			if s := os.Getenv(o.Environment); s != "" {
				const base, bitSize = 10, 64
				v, err := strconv.ParseInt(s, base, bitSize)
				if err != nil {
					return errorz.Errorf("%s: %w", o.Environment, err)
				}
				o.value = &v
			}
		case *Float64Option:
			if s := os.Getenv(o.Environment); s != "" {
				const bitSize = 64
				v, err := strconv.ParseFloat(s, bitSize)
				if err != nil {
					return errorz.Errorf("%s: %w", o.Environment, err)
				}
				o.value = &v
			}
		case *HelpOption:
			// do nothing
		default:
			return errorz.Errorf("%s: %w", o.GetName(), ErrInvalidOptionType)
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.loadEnvironments(); err != nil {
			return err
		}
	}
	return nil
}
