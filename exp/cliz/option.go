package cliz

import (
	"strings"
)

type (
	// Option is the interface for the option.
	Option interface {
		// GetName returns the name of the option.
		GetName() string
		// GetAliases returns the alias names of the option.
		GetAliases() []string
		// GetEnvironment returns the environment variable name of the option.
		GetEnvironment() string
		// GetDescription returns the description of the option.
		GetDescription() string
		// GetDefault returns the default value of the option.
		GetDefault() interface{}
		// IsRequired returns the required flag of the option.
		IsRequired() bool
		// IsZero returns whether the option has a value.
		IsZero() bool
	}
)

// --long or -s
func argIsHyphenOption(o Option, arg string) bool {
	return arg == longOptionPrefix+o.GetName() ||
		func() bool {
			for _, alias := range o.GetAliases() {
				if arg == shortOptionPrefix+alias || arg == longOptionPrefix+alias {
					return true
				}
			}
			return false
		}()
}

// --long=value or -s=value
func argIsHyphenOptionEqual(o Option, arg string) bool {
	return strings.HasPrefix(arg, longOptionPrefix+o.GetName()+"=") ||
		func() bool {
			for _, alias := range o.GetAliases() {
				if strings.HasPrefix(arg, shortOptionPrefix+alias+"=") || strings.HasPrefix(arg, longOptionPrefix+alias+"=") {
					return true
				}
			}
			return false
		}()
}

func extractValueFromHyphenOptionEqual(arg string) string {
	return strings.Join(strings.Split(arg, "=")[1:], "=")
}

func hasNoOptionValue(args []string, i int) bool {
	lastIndex := len(args) - 1
	return i+1 > lastIndex
}
