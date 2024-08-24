package cliz

import (
	"os"
	"strconv"

	"github.com/hakadoriya/z.go/cliz/clierrz"
	"github.com/hakadoriya/z.go/errorz"
)

//nolint:gocognit,cyclop
func (c *Command) loadEnvironments() error {
	for _, opt := range c.Options {
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
		default:
			return errorz.Errorf("%s: %w", o.GetName(), clierrz.ErrInvalidOptionType)
		}
	}

	for _, subcmd := range c.SubCommands {
		if err := subcmd.loadEnvironments(); err != nil {
			return err
		}
	}
	return nil
}
