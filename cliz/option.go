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
		// GetEnv returns the environment variable name of the option.
		GetEnv() string
		// GetDefault returns the default value of the option.
		GetDefault() interface{}
		// IsRequired returns the required flag of the option.
		IsRequired() bool
		// IsZero returns whether the option has a value.
		IsZero() bool
		// IsHidden returns the hidden flag of the option.
		IsHidden() bool
		// GetDescription returns the description of the option.
		GetDescription() string
	}
)

// --long or -s
func argIsHyphenOption(o Option, osArg string) bool {
	return osArg == longOptionPrefix+o.GetName() ||
		func() bool {
			for _, alias := range o.GetAliases() {
				if osArg == shortOptionPrefix+alias || osArg == longOptionPrefix+alias {
					return true
				}
			}
			return false
		}()
}

// --long=value or -s=value
func argIsHyphenOptionEqual(o Option, osArg string) bool {
	return strings.HasPrefix(osArg, longOptionPrefix+o.GetName()+"=") ||
		func() bool {
			for _, alias := range o.GetAliases() {
				if strings.HasPrefix(osArg, shortOptionPrefix+alias+"=") || strings.HasPrefix(osArg, longOptionPrefix+alias+"=") {
					return true
				}
			}
			return false
		}()
}

func extractValueFromHyphenOptionEqual(osArg string) string {
	return strings.Join(strings.Split(osArg, "=")[1:], "=")
}

func hasNoOptionValue(osArgs []string, i int) bool {
	lastIndex := len(osArgs) - 1
	return i+1 > lastIndex
}
