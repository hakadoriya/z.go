package slogz

import "log/slog"

type HandlerOption interface {
	apply(s *slogHandler)
}

type withHandlerOptions struct{ o *slog.HandlerOptions }

func (o withHandlerOptions) apply(s *slogHandler) { s.slogHandlerOptions = o.o }

func WithHandlerOptions(o *slog.HandlerOptions) HandlerOption {
	return withHandlerOptions{o: o}
}

type withErrorVerbose struct{ errorVerbose bool }

func (o withErrorVerbose) apply(s *slogHandler) { s.errorVerbose = o.errorVerbose }

func WithErrorVerbose(verbose bool) HandlerOption {
	return withErrorVerbose{errorVerbose: verbose}
}

type withErrorVerboseKeySuffix struct{ suffix string }

func (o withErrorVerboseKeySuffix) apply(s *slogHandler) { s.errorVerboseKeySuffix = o.suffix }

func WithErrorVerboseKeySuffix(suffix string) HandlerOption {
	return withErrorVerboseKeySuffix{suffix: suffix}
}
