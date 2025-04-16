package errorz

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
)

type (
	ErrorfOption interface {
		apply(c *errorfConfig)
	}
	errorfConfig struct {
		addCallerSkip int
	}
)

type withErrorfOptionAddCallerSkip struct{ addCallerSkip int }

func (o withErrorfOptionAddCallerSkip) apply(c *errorfConfig) {
	c.addCallerSkip += o.addCallerSkip
}

// WithErrorfOptionAddCallerSkip returns an ErrorfOption that sets the number of stack frames to skip.
func WithErrorfOptionAddCallerSkip(skip int) ErrorfOption {
	return withErrorfOptionAddCallerSkip{skip}
}

// NewErrorf returns an Errorf function with stack trace capabilities similar to xerrors.Errorf.
//
// Differences from xerrors.Errorf:
// xerrors.Errorf cannot set caller skip, so this function was created.
// NewErrorf allows setting caller skip via the WithAddCallerSkip option.
//
// For example, if you want to add the call of a utility function that receives an error to the error stack trace,
// xerrors.Errorf cannot set caller skip, so you cannot add the caller of that utility function to the error stack trace.
// NewErrorf allows you to create a new Errorf function with added caller skip individually.
// By using this individually created Errorf function within the utility function,
// you can add the caller of the utility function to the error stack trace.
func NewErrorf(opts ...ErrorfOption) func(format string, a ...interface{}) error {
	c := new(errorfConfig)

	for _, opt := range opts {
		opt.apply(c)
	}

	return newErrorf(c)
}

const (
	indent4 = "    "
	ln      = "\n"
)

//nolint:cyclop
func newErrorf(c *errorfConfig) func(format string, a ...interface{}) error {
	return func(format string, a ...interface{}) error {
		const (
			suffixS      = ": %s"
			suffixV      = ": %v"
			suffixPlusV  = ": %+v"
			suffixSharpV = ": %#v"
			suffixW      = ": %w"
		)
		var (
			hasSuffixS      = strings.HasSuffix(format, suffixS)
			hasSuffixV      = strings.HasSuffix(format, suffixV)
			hasSuffixPlusV  = strings.HasSuffix(format, suffixPlusV)
			hasSuffixSharpV = strings.HasSuffix(format, suffixSharpV)
			hasSuffixW      = strings.HasSuffix(format, suffixW)
		)

		if !hasSuffixS && !hasSuffixV && !hasSuffixPlusV && !hasSuffixSharpV && !hasSuffixW {
			return fmt.Errorf(format, a...) //nolint:err113
		}

		prefix := format[:len(format)-len(suffixW)]
		suffix := format[len(format)-len(suffixW):]
		head := a[:len(a)-1]
		tail := a[len(a)-1]

		var e wrapError
		runtime.Callers(1+c.addCallerSkip, e.frame[:])
		e.msg = fmt.Sprintf(prefix, head...)
		switch err := tail.(type) {
		case formatter:
			e.err = err
		case error:
			switch {
			case hasSuffixS:
				e.err = fmt.Errorf("%s", err) //nolint:errorlint,err113 // for compatibility with xerrors.Errorf
			case hasSuffixV || hasSuffixPlusV || hasSuffixSharpV:
				e.err = fmt.Errorf("%v", err) //nolint:errorlint,err113 // for compatibility with xerrors.Errorf
			// case hasSuffixPlusV: // FIXME: support %+v
			// 	e.err = fmt.Errorf("%+v", err) //nolint:errorlint,err113 // for compatibility with xerrors.Errorf
			// case hasSuffixSharpV: // FIXME: support %#v
			// 	e.err = fmt.Errorf("%+v", err) //nolint:errorlint,err113 // for compatibility with xerrors.Errorf
			case hasSuffixW:
				e.err = err
			}
		default:
			e.msg += fmt.Sprintf(suffix, tail)
			e.err = nil
		}

		return &e
	}
}

//nolint:gochecknoglobals
var errorf = NewErrorf(WithErrorfOptionAddCallerSkip(1))

// Errorf is a function similar to xerrors.Errorf.
// It uses the wrapError type, which satisfies the error interface, to retain the stack trace.
// Additionally, it implements the fmt.Formatter interface, allowing the stack trace to be displayed using fmt.Printf and similar functions.
func Errorf(format string, a ...interface{}) error {
	return errorf(format, a...)
}

type wrapError struct {
	msg   string
	err   error
	frame [3]uintptr // See: https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
}

var (
	_ error                       = (*wrapError)(nil)
	_ formatter                   = (*wrapError)(nil)
	_ fmt.Formatter               = (*wrapError)(nil)
	_ fmt.GoStringer              = (*wrapError)(nil)
	_ interface{ Unwrap() error } = (*wrapError)(nil)
)

type formatter interface {
	error
	format(s fmt.State, verb rune)
	Unwrap() error
}

func (e *wrapError) Error() string {
	return fmt.Sprint(e)
}

func (e *wrapError) Format(s fmt.State, verb rune) {
	var err error = e
loop:
	for {
		switch fe := err.(type) { //nolint:errorlint
		case formatter:
			fe.format(s, verb)
			err = fe.Unwrap()
		case fmt.Formatter:
			fe.Format(s, verb)
			break loop
		default:
			_, _ = fmt.Fprintf(s, fmt.FormatString(s, verb), fe)
			break loop
		}
		if err == nil {
			break loop
		}
	}
}

func (e *wrapError) GoString() string {
	valType := reflect.TypeOf(*e)
	val := reflect.ValueOf(*e)
	elems := make([]string, valType.NumField())
	for i := range valType.NumField() {
		elems[i] = fmt.Sprintf("%s:%#v", valType.Field(i).Name, val.Field(i))
	}
	return fmt.Sprintf("&%s{%s}", valType, strings.Join(elems, ", "))
}

func (e *wrapError) Unwrap() error {
	return e.err
}

// FormatError is intended to be used as follows:
//
//	func (e *customError) Format(s fmt.State, verb rune) {
//		errorz.FormatError(s, verb, e.Unwrap())
//	}
func FormatError(s fmt.State, verb rune, err error) {
	if formatter, ok := err.(fmt.Formatter); ok {
		formatter.Format(s, verb)
		return
	}

	_, _ = fmt.Fprintf(s, fmt.FormatString(s, verb), err)
}

func (e *wrapError) writeCallers(w io.Writer) {
	frames := runtime.CallersFrames(e.frame[:])
	if _, ok := frames.Next(); !ok {
		return
	}
	target, ok := frames.Next()
	if !ok {
		return
	}

	if target.Function != "" {
		fmt.Fprintf(w, ":"+ln+indent4+"%s", target.Function)
		// NOTE:
		//              ^^^^^^^^^^^^^^^^^
		//              means a part of stacktrace:
		//
		// funcA:\n
		//      ^^^
		//     github.com/org/repo/pkg.funcA
		// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
		if target.File != "" {
			fmt.Fprintf(w, ln+indent4+indent4+"%s:%d", target.File, target.Line)
			// NOTE:
			//             ^^^^^^^^^^^^^^^^^^^^^^^^^
			//             means a part of stacktrace:
			//
			//     github.com/org/repo/pkg.funcA\n
			//                                  ^^
			//         github.com/org/repo/pkg.go:123
			// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
		}
	}
}

func (e *wrapError) format(s fmt.State, verb rune) {
	var withStacktrace bool
Verb:
	switch verb {
	// FormatError() will not be called with the 'w' verb.
	// case 'w':
	case 'v':
		switch {
		case s.Flag('#'):
			_, _ = io.WriteString(s, e.GoString())
			return
		case s.Flag('+'):
			withStacktrace = true
			break Verb
		}
	default:
	}

	_, _ = io.WriteString(s, e.msg)
	if withStacktrace {
		e.writeCallers(s)
		if e.err != nil {
			_, _ = io.WriteString(s, ln+"  - ")
			// NOTE:
			//                        ^^^^^^
			//                        means a part of stacktrace:
			//
			//         github.com/org/repo/pkg.go:123\n
			//                                       ^^
			//   - funcB:
			// ^^^^
		}
	} else { //nolint:gocritic
		if e.err != nil {
			_, _ = io.WriteString(s, ": ")
			// NOTE:
			//                        ^^
			//                        means a part of error output:
			// funcA: funcB: funcC: error
			//      ^^     ^^     ^^
		}
	}
}
