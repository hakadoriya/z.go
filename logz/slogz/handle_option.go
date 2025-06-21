package slogz

import "log/slog"

type HandlerOption interface {
	apply(s *slogJSONHandler)
}

type handlerOptionFunc func(s *slogJSONHandler)

func (f handlerOptionFunc) apply(s *slogJSONHandler) { f(s) }

func WithHandlerOptionHandlerOptions(o *slog.HandlerOptions) HandlerOption {
	return handlerOptionFunc(func(s *slogJSONHandler) { s.slogHandlerOptions = o })
}

func WithHandlerOptionAddSource(addSource bool) HandlerOption {
	return handlerOptionFunc(func(s *slogJSONHandler) { s.slogHandlerOptions.AddSource = addSource })
}

func WithHandlerOptionAddCallerSkip(skip int) HandlerOption {
	return handlerOptionFunc(func(s *slogJSONHandler) { s.addCallerSkip += skip })
}

func WithHandlerOptionAddTimestamp(addTimestamp bool) HandlerOption {
	return handlerOptionFunc(func(s *slogJSONHandler) { s.addTimestamp = addTimestamp })
}

func WithHandlerOptionAddAttrs(attrs ...slog.Attr) HandlerOption {
	return handlerOptionFunc(func(s *slogJSONHandler) { s.addAttrs = attrs })
}

func WithHandlerOptionErrorVerbose(verbose bool) HandlerOption {
	return handlerOptionFunc(func(s *slogJSONHandler) { s.errorVerbose = verbose })
}

func WithHandlerOptionErrorVerboseKeySuffix(suffix string) HandlerOption {
	return handlerOptionFunc(func(s *slogJSONHandler) { s.errorVerboseKeySuffix = suffix })
}
