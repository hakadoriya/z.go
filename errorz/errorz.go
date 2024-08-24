package errorz

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

func Contains(err error, substr string) bool {
	return err != nil && strings.Contains(err.Error(), substr)
}

func HasPrefix(err error, prefix string) bool {
	return err != nil && strings.HasPrefix(err.Error(), prefix)
}

func HasSuffix(err error, suffix string) bool {
	return err != nil && strings.HasSuffix(err.Error(), suffix)
}

func MatchRegex(err error, re *regexp.Regexp) bool {
	return err != nil && re.MatchString(err.Error())
}

func IsNetTimeout(err error) bool {
	if netErr := (net.Error)(nil); errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return false
}
