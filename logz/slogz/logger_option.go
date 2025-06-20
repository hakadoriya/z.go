package slogz

type slogJSONLoggerConfig struct {
	handlerOptions []HandlerOption
}

type LoggerOption interface {
	apply(c *slogJSONLoggerConfig)
}

type loggerOptionFunc func(c *slogJSONLoggerConfig)

func (f loggerOptionFunc) apply(c *slogJSONLoggerConfig) { f(c) }

func WithLoggerOptionHandlerOption(opts ...HandlerOption) LoggerOption {
	return loggerOptionFunc(func(c *slogJSONLoggerConfig) { c.handlerOptions = append(c.handlerOptions, opts...) })
}

func WithLoggerOptionAddCallerSkip(skip int) LoggerOption {
	return loggerOptionFunc(func(c *slogJSONLoggerConfig) {
		c.handlerOptions = append(c.handlerOptions, WithHandlerOptionAddCallerSkip(skip))
	})
}
