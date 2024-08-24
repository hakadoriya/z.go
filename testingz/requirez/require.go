package requirez

import (
	"regexp"
	"testing"

	"github.com/hakadoriya/z.go/testingz/internal"
)

// NoError asserts that err is nil.
func NoError(tb testing.TB, err error, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.NoError(tb, tb.Fatal, err, formatAndArgs...)
}

// Error asserts that err is not nil.
func Error(tb testing.TB, err error, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.Error(tb, tb.Fatal, err, formatAndArgs...)
}

// ErrorIs asserts that err is target.
func ErrorIs(tb testing.TB, err, target error, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.ErrorIs(tb, tb.Fatal, err, target, formatAndArgs...)
}

// ErrorContains asserts that err contains substr.
func ErrorContains(tb testing.TB, err error, substr string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.ErrorContains(tb, tb.Fatal, err, substr, formatAndArgs...)
}

// ErrorMatchRegex asserts that err matches pattern.
func ErrorMatchRegex(tb testing.TB, err error, re *regexp.Regexp, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.ErrorMatchRegex(tb, tb.Fatal, err, re, formatAndArgs...)
}

func StringHasPrefix(tb testing.TB, s, prefix string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringHasPrefix(tb, tb.Fatal, s, prefix, formatAndArgs...)
}

func StringHasSuffix(tb testing.TB, s, suffix string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringHasSuffix(tb, tb.Fatal, s, suffix, formatAndArgs...)
}

// StringContains asserts that s contains substr.
func StringContains(tb testing.TB, s, substr string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringContains(tb, tb.Fatal, s, substr, formatAndArgs...)
}

// StringMatchRegex asserts that s matches pattern.
func StringMatchRegex(tb testing.TB, s string, re *regexp.Regexp, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringMatchRegex(tb, tb.Fatal, s, re, formatAndArgs...)
}

// True asserts that value is true.
func True(tb testing.TB, value bool, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.True(tb, tb.Fatal, value, formatAndArgs...)
}

// False asserts that value is false.
func False(tb testing.TB, value bool, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.False(tb, tb.Fatal, value, formatAndArgs...)
}

// Equal asserts that expected and actual are deeply equal.
func Equal(tb testing.TB, expected, actual interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.Equal(tb, tb.Fatal, expected, actual, formatAndArgs...)
}

// NotEqual asserts that expected and actual are not deeply equal.
func NotEqual(tb testing.TB, expected, actual interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.NotEqual(tb, tb.Fatal, expected, actual, formatAndArgs...)
}

// Nil asserts that value is nil.
func Nil(tb testing.TB, value interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.Nil(tb, tb.Fatal, value, formatAndArgs...)
}

// NotNil asserts that value is not nil.
func NotNil(tb testing.TB, value interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.NotNil(tb, tb.Fatal, value, formatAndArgs...)
}
