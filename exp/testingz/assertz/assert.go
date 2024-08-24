package assertz

import (
	"regexp"
	"testing"

	"github.com/hakadoriya/z.go/exp/testingz/internal"
)

// NoError asserts that err is nil.
func NoError(tb testing.TB, err error, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.NoError(tb, tb.Error, err, msgAndArgs...)
}

// Error asserts that err is not nil.
func Error(tb testing.TB, err error, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.Error(tb, tb.Error, err, msgAndArgs...)
}

// ErrorIs asserts that err is target.
func ErrorIs(tb testing.TB, err, target error, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.ErrorIs(tb, tb.Error, err, target, msgAndArgs...)
}

// ErrorContains asserts that err contains substr.
func ErrorContains(tb testing.TB, err error, substr string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.ErrorContains(tb, tb.Error, err, substr, msgAndArgs...)
}

// ErrorMatchRegex asserts that err matches pattern.
func ErrorMatchRegex(tb testing.TB, err error, re *regexp.Regexp, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.ErrorMatchRegex(tb, tb.Error, err, re, msgAndArgs...)
}

func StringHasPrefix(tb testing.TB, s, prefix string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringHasPrefix(tb, tb.Error, s, prefix, msgAndArgs...)
}

func StringHasSuffix(tb testing.TB, s, suffix string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringHasSuffix(tb, tb.Error, s, suffix, msgAndArgs...)
}

// StringContains asserts that s contains substr.
func StringContains(tb testing.TB, s, substr string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringContains(tb, tb.Error, s, substr, msgAndArgs...)
}

// StringMatchRegex asserts that s matches pattern.
func StringMatchRegex(tb testing.TB, s string, re *regexp.Regexp, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.StringMatchRegex(tb, tb.Error, s, re, msgAndArgs...)
}

// True asserts that value is true.
func True(tb testing.TB, value bool, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.True(tb, tb.Error, value, msgAndArgs...)
}

// False asserts that value is false.
func False(tb testing.TB, value bool, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.False(tb, tb.Error, value, msgAndArgs...)
}

// Equal asserts that expected and actual are deeply equal.
func Equal(tb testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.Equal(tb, tb.Error, expected, actual, msgAndArgs...)
}

// NotEqual asserts that expected and actual are not deeply equal.
func NotEqual(tb testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.NotEqual(tb, tb.Error, expected, actual, msgAndArgs...)
}

// Nil asserts that value is nil.
func Nil(tb testing.TB, value interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.Nil(tb, tb.Error, value, msgAndArgs...)
}

// NotNil asserts that value is not nil.
func NotNil(tb testing.TB, value interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	return internal.NotNil(tb, tb.Error, value, msgAndArgs...)
}
