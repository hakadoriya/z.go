package runtimez

import (
	"path"
	"runtime"
)

type FuncNameOption interface {
	apply(c *funcNameConfig)
}

type funcNameConfig struct {
	skip int
}

type withFuncNameOptionAddCallerSkip struct {
	skip int
}

func (w *withFuncNameOptionAddCallerSkip) apply(c *funcNameConfig) {
	c.skip += w.skip
}

func WithFuncNameOptionAddCallerSkip(skip int) FuncNameOption {
	return &withFuncNameOptionAddCallerSkip{skip: skip}
}

// FuncName returns the name of the function that executed this function.
// The skip parameter specifies how many levels up the call stack to go to retrieve the function name.
//
// Example:
//
//	func main() {
//		fmt.Println(wrapFunc()) // Output -> main.wrapFunc
//	}
//
//	func wrapFunc() string {
//		return FuncName()
//	}
func FuncName(opts ...FuncNameOption) (funcName string) {
	return path.Base(FullFuncName(append(opts, WithFuncNameOptionAddCallerSkip(1))...))
}

// The output of FullFuncName is usually too long, so typically use FuncName instead.
//
// Example:
//
//	func main() {
//		fmt.Println(wrapFunc()) // Output -> github.com/hakadoriya/z.go/main.wrapFunc
//	}
//
//	func wrapFunc() string {
//		return FuncName()
//	}
func FullFuncName(opts ...FuncNameOption) (funcName string) {
	cfg := &funcNameConfig{
		skip: 1,
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	pc, _, _, _ := runtime.Caller(cfg.skip) //nolint:dogsled
	return runtime.FuncForPC(pc).Name()
}
