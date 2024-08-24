package internal

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/exp/diffz/simplediffz"
	stringz "github.com/hakadoriya/z.go/stringz"
)

func NoError(tb testing.TB, printFunc func(args ...any), err error, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if err != nil {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err != nil: %+v", tb.Name(), err)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s: %+v", tb.Name(), fmt.Sprint(msgAndArgs...), err)
		default:
			msg = fmt.Sprintf("❌: %s: %s: %+v", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...), err)
		}

		printFunc(msg)

		return false
	}
	return true
}

func Error(tb testing.TB, printFunc func(args ...any), err error, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if err == nil {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err == nil", tb.Name())
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s", tb.Name(), fmt.Sprint(msgAndArgs...))
		default:
			msg = fmt.Sprintf("❌: %s: %s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
		}

		printFunc(msg)

		return false
	}
	return true
}

func ErrorIs(tb testing.TB, printFunc func(args ...any), err error, target error, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !errors.Is(err, target) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err != target:\n--- TARGET\n+++ ERROR\n%s\n%s",
				tb.Name(),
				stringz.AddPrefix("-", fmt.Sprintf("%v", target), "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- TARGET\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprint(msgAndArgs...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", target), "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- TARGET\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", target), "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		}

		printFunc(msg)

		return false
	}
	return true
}

func ErrorContains(tb testing.TB, printFunc func(args ...any), err error, substr string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !errorz.Contains(err, substr) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err not contains substr:\n--- SUBSTR\n+++ ERROR\n%s\n%s",
				tb.Name(),
				stringz.AddPrefix("-", substr, "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- SUBSTR\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprint(msgAndArgs...),
				stringz.AddPrefix("-", substr, "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- SUBSTR\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...),
				stringz.AddPrefix("-", substr, "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func ErrorMatchRegex(tb testing.TB, printFunc func(args ...any), err error, re *regexp.Regexp, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !errorz.MatchRegex(err, re) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err not match regex:\n--- REGEX\n+++ ERROR\n%s\n%s",
				tb.Name(),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- REGEX\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprint(msgAndArgs...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- REGEX\n+++ ERROR\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), "\n"),
				stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringHasPrefix(tb testing.TB, printFunc func(args ...any), s string, prefix string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !strings.HasPrefix(s, prefix) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err not has prefix:\n--- PREFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				stringz.AddPrefix("-", prefix, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- PREFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(msgAndArgs...),
				stringz.AddPrefix("-", prefix, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- PREFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...),
				stringz.AddPrefix("-", prefix, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringHasSuffix(tb testing.TB, printFunc func(args ...any), s string, suffix string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !strings.HasSuffix(s, suffix) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err not has suffix:\n--- SUFFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				stringz.AddPrefix("-", suffix, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- SUFFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(msgAndArgs...),
				stringz.AddPrefix("-", suffix, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- SUFFIX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...),
				stringz.AddPrefix("-", suffix, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringContains(tb testing.TB, printFunc func(args ...any), s string, substr string, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !strings.Contains(s, substr) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err not contains substr:\n--- SUBSTR\n+++ STRING\n%s\n%s",
				tb.Name(),
				stringz.AddPrefix("-", substr, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- SUBSTR\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(msgAndArgs...),
				stringz.AddPrefix("-", substr, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- SUBSTR\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...),
				stringz.AddPrefix("-", substr, "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func StringMatchRegex(tb testing.TB, printFunc func(args ...any), s string, re *regexp.Regexp, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !re.MatchString(s) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: err not match regex:\n--- REGEX\n+++ STRING\n%s\n%s",
				tb.Name(),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- REGEX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprint(msgAndArgs...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- REGEX\n+++ STRING\n%s\n%s",
				tb.Name(),
				fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...),
				stringz.AddPrefix("-", fmt.Sprintf("%v", re), "\n"),
				stringz.AddPrefix("+", s, "\n"),
			)
		}

		printFunc(msg)

		return false
	}

	return true
}

func True(tb testing.TB, printFunc func(args ...any), value bool, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !value {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: value == false", tb.Name())
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s", tb.Name(), fmt.Sprint(msgAndArgs...))
		default:
			msg = fmt.Sprintf("❌: %s: %s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
		}

		printFunc(msg)

		return false
	}
	return true
}

func False(tb testing.TB, printFunc func(args ...any), value bool, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if value {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: value == true", tb.Name())
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s", tb.Name(), fmt.Sprint(msgAndArgs...))
		default:
			msg = fmt.Sprintf("❌: %s: %s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
		}

		printFunc(msg)

		return false
	}
	return true
}

func Equal(tb testing.TB, printFunc func(args ...any), expected, actual interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if !reflect.DeepEqual(expected, actual) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: expected != actual:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", expected, expected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)))
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprint(msgAndArgs...), simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", expected, expected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)))
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...), simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", expected, expected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)))
		}

		printFunc(msg)

		return false
	}
	return true
}

func NotEqual(tb testing.TB, printFunc func(args ...any), unexpected, actual interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()

	if reflect.DeepEqual(unexpected, actual) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: expected == actual:\n--- UNEXPECTED\n+++ ACTUAL\n%s", tb.Name(), simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", unexpected, unexpected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)))
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- UNEXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprint(msgAndArgs...), simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", unexpected, unexpected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)))
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- UNEXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...), simplediffz.Diff(fmt.Sprintf("TYPE:%T\n%+v", unexpected, unexpected), fmt.Sprintf("TYPE:%T\n%+v", actual, actual)))
		}

		printFunc(msg)

		return false
	}
	return true
}

func Nil(tb testing.TB, printFunc func(args ...any), value interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()
	defer func() {
		tb.Helper()
		if r := recover(); r != nil {
			var msg string
			switch {
			case len(msgAndArgs) == 0:
				msg = fmt.Sprintf("❌: %s: value != nil:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)))
			case len(msgAndArgs) == 1:
				msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprint(msgAndArgs...), simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)))
			default:
				msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...), simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)))
			}

			printFunc(msg)

			success = false
		}
	}()

	if !(value == nil || reflect.ValueOf(value).IsNil()) {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: value != nil:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)))
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprint(msgAndArgs...), simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)))
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...), simplediffz.Diff(fmt.Sprintf("%+v", nil), fmt.Sprintf("%+v", value)))
		}

		printFunc(msg)

		return false
	}
	return true
}

func NotNil(tb testing.TB, printFunc func(args ...any), value interface{}, msgAndArgs ...interface{}) (success bool) {
	tb.Helper()
	defer func() {
		tb.Helper()
		if r := recover(); r != nil {
			var msg string
			switch {
			case len(msgAndArgs) == 0:
				msg = fmt.Sprintf("❌: %s: value == nil:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)))
			case len(msgAndArgs) == 1:
				msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprint(msgAndArgs...), simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)))
			default:
				msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...), simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)))
			}

			printFunc(msg)

			success = false
		}
	}()

	if value == nil || reflect.ValueOf(value).IsNil() {
		var msg string
		switch {
		case len(msgAndArgs) == 0:
			msg = fmt.Sprintf("❌: %s: value == nil:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)))
		case len(msgAndArgs) == 1:
			msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprint(msgAndArgs...), simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)))
		default:
			msg = fmt.Sprintf("❌: %s: %s:\n--- EXPECTED\n+++ ACTUAL\n%s", tb.Name(), fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...), simplediffz.Diff("NOT <nil>", fmt.Sprintf("%+v", value)))
		}

		printFunc(msg)

		return false
	}
	return true
}
