package panicz

import "errors"

func Panic(err error, opts ...PanicOption) {
	if err == nil {
		return
	}

	config := new(panicConfig)
	for _, opt := range opts {
		opt.apply(config)
	}

	for _, ignore := range config.ignores {
		if errors.Is(err, ignore) {
			return
		}
	}

	panic(err)
}

type (
	PanicOption     interface{ apply(config *panicConfig) }
	panicOptionFunc func(config *panicConfig)
	panicConfig     struct {
		ignores []error
	}
)

func (f panicOptionFunc) apply(config *panicConfig) { f(config) }

func WithPanicOptionIgnoreErrors(ignores ...error) PanicOption {
	return panicOptionFunc(func(config *panicConfig) {
		config.ignores = ignores
	})
}
