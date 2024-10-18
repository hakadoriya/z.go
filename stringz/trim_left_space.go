package stringz

import (
	"strings"
	"unicode"
)

// TrimLeftSpace trims the leading space characters from a string.
func TrimLeftSpace(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

// TrimRightSpace trims the trailing space characters from a string.
func TrimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}
