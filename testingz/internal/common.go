package internal

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hakadoriya/z.go/diffz/simplediffz"
	"github.com/hakadoriya/z.go/errorz"
	stringz "github.com/hakadoriya/z.go/stringz"
	"github.com/hakadoriya/z.go/testingz"
)

func NoError(tb testing.TB, printFunc func(args ...any), err error, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if err != nil {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s: %+v",
				tb.Name(),
				"err != nil",
				err,
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s: %+v",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				err,
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s: %+v",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				err,
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

func Error(tb testing.TB, printFunc func(args ...any), err error, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if err == nil {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				"err == nil",
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

func ErrorIs(tb testing.TB, printFunc func(args ...any), err error, target error, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !errors.Is(err, target) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- TARGET\n+++ ERROR\n%s\n%s",
				tb.Name(),
				"err != target",
				stringz.AddPrefix("-", fmt.Sprintf("%v", target), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- TARGET\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", target), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- TARGET\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", target), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

func ErrorContains(tb testing.TB, printFunc func(args ...any), err error, substr string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !errorz.Contains(err, substr) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUBSTR\n+++ ERROR\n%s\n%s",
				tb.Name(),
				"err not contains substr",
				stringz.AddPrefix("-", substr, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUBSTR\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				stringz.AddPrefix("-", substr, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUBSTR\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				stringz.AddPrefix("-", substr, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func ErrorMatchRegex(tb testing.TB, printFunc func(args ...any), err error, re *regexp.Regexp, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !errorz.MatchRegex(err, re) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- REGEX\n+++ ERROR\n%s\n%s",
				tb.Name(),
				"err not match regex",
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- REGEX\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- REGEX\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), testingz.DefaultLineSeparator),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringHasPrefix(tb testing.TB, printFunc func(args ...any), s string, prefix string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !strings.HasPrefix(s, prefix) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- PREFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				"err not has prefix",
				stringz.AddPrefix("-", prefix, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- PREFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				stringz.AddPrefix("-", prefix, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- PREFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				stringz.AddPrefix("-", prefix, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringHasSuffix(tb testing.TB, printFunc func(args ...any), s string, suffix string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !strings.HasSuffix(s, suffix) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUFFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				"err not has suffix",
				stringz.AddPrefix("-", suffix, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUFFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				stringz.AddPrefix("-", suffix, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUFFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				stringz.AddPrefix("-", suffix, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringContains(tb testing.TB, printFunc func(args ...any), s string, substr string, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !strings.Contains(s, substr) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUBSTR\n+++ STRING\n%s\n%s",
				tb.Name(),
				"err not contains substr",
				stringz.AddPrefix("-", substr, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUBSTR\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				stringz.AddPrefix("-", substr, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- SUBSTR\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				stringz.AddPrefix("-", substr, testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringMatchRegex(tb testing.TB, printFunc func(args ...any), s string, re *regexp.Regexp, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !re.MatchString(s) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- REGEX\n+++ STRING\n%s\n%s",
				tb.Name(),
				"err not match regex",
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- REGEX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- REGEX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), testingz.DefaultLineSeparator),
				stringz.AddPrefix("+", s, testingz.DefaultLineSeparator),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func True(tb testing.TB, printFunc func(args ...any), value bool, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !value {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				"value == false",
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

func False(tb testing.TB, printFunc func(args ...any), value bool, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if value {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				"value == true",
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

func Equal(tb testing.TB, printFunc func(args ...any), expected, actual interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !reflect.DeepEqual(expected, actual) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				"expected != actual",
				simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", expected, expected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", expected, expected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", expected, expected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

func NotEqual(tb testing.TB, printFunc func(args ...any), unexpected, actual interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if reflect.DeepEqual(unexpected, actual) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- UNEXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				"expected == actual",
				simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", unexpected, unexpected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- UNEXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", unexpected, unexpected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- UNEXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", unexpected, unexpected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

//nolint:funlen
func Nil(tb testing.TB, printFunc func(args ...any), value interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()
	defer func() {
		tb.Helper()
		if r := recover(); r != nil {
			var msg string
			switch {
			case len(formatAndArgs) == 0:
				msg = fmt.Sprintf(
					testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
					tb.Name(),
					"value != nil",
					simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)),
				)
			case len(formatAndArgs) == 1:
				msg = fmt.Sprintf(
					testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
					tb.Name(),
					fmt.Sprint(formatAndArgs...),
					simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)),
				)
			default:
				msg = fmt.Sprintf(
					testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
					tb.Name(),
					fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
					simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)),
				)
			}

			printFunc(msg)

			success = false
		}
	}()

	if !(value == nil || reflect.ValueOf(value).IsNil()) {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				"value != nil",
				simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

//nolint:funlen
func NotNil(tb testing.TB, printFunc func(args ...any), value interface{}, formatAndArgs ...interface{}) (success bool) {
	tb.Helper()
	defer func() {
		tb.Helper()
		if r := recover(); r != nil {
			var msg string
			switch {
			case len(formatAndArgs) == 0:
				msg = fmt.Sprintf(
					testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
					tb.Name(),
					"value == nil",
					simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)),
				)
			case len(formatAndArgs) == 1:
				msg = fmt.Sprintf(
					testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
					tb.Name(),
					fmt.Sprint(formatAndArgs...),
					simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)),
				)
			default:
				msg = fmt.Sprintf(
					testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
					tb.Name(),
					fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
					simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)),
				)
			}

			printFunc(msg)

			success = false
		}
	}()

	if value == nil || reflect.ValueOf(value).IsNil() {
		var msg string
		switch {
		case len(formatAndArgs) == 0:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				"value == nil",
				simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)),
			)
		case len(formatAndArgs) == 1:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprint(formatAndArgs...),
				simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)),
			)
		default:
			msg = fmt.Sprintf(
				testingz.DefaultFailMarker+"%s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(formatAndArgs[0]), formatAndArgs[1:]...),
				simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}
