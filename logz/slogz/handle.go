package slogz

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"time"
)

var _ slog.Handler = (*slogJSONHandler)(nil)

type slogJSONHandler struct {
	slogHandlerOptions    *slog.HandlerOptions
	w                     io.Writer
	addCallerSkip         int
	addTimestamp          bool
	addAttrs              []slog.Attr
	errorVerbose          bool
	errorVerboseKeySuffix string
	slogHandler           slog.Handler
}

func NewHandler(w io.Writer, level slog.Leveler, opts ...HandlerOption) slog.Handler {
	return newHandler(w, level, opts...)
}

func newHandler(w io.Writer, level slog.Leveler, opts ...HandlerOption) *slogJSONHandler {
	const defaultCallerSkip = 4
	s := &slogJSONHandler{
		slogHandlerOptions: &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: ReplaceAttr,
		},
		w:                     w,
		addCallerSkip:         defaultCallerSkip,
		addTimestamp:          true,
		addAttrs:              nil,
		errorVerbose:          true,
		errorVerboseKeySuffix: "Verbose",
		slogHandler:           nil,
	}

	for _, o := range opts {
		o.apply(s)
	}

	s.slogHandler = slog.NewJSONHandler(s.w, s.slogHandlerOptions)

	return s
}

func RenewHandler(ctx context.Context, handler slog.Handler, opts ...HandlerOption) slog.Handler {
	source, ok := handler.(*slogJSONHandler)
	if !ok {
		const defaultCallerSkip = 1
		var pcs [1]uintptr
		runtime.Callers(defaultCallerSkip, pcs[:])
		err := fmt.Errorf("RenewHandler: handler=%T: %w", handler, ErrHandlerIsNotSlogJSONHandler)
		r := slog.NewRecord(time.Now(), slog.LevelWarn, err.Error(), pcs[0])
		r.AddAttrs(Error(err))
		_ = handler.Handle(ctx, r)
		source = newHandler(DefaultWriter, DefaultLevel)
	}

	s := source.clone()
	for _, o := range opts {
		o.apply(s)
	}

	s.slogHandler = slog.NewJSONHandler(s.w, s.slogHandlerOptions)

	return s
}

func (s *slogJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.slogHandler.Enabled(ctx, level)
}

func (s *slogJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		// Attr Value type switch
		switch v := a.Value.Any().(type) {
		case error:
			// If errorVerbose is set, add verbose error to the record.
			if s.errorVerbose {
				attrs = append(attrs, slog.String(a.Key+s.errorVerboseKeySuffix, fmt.Sprintf("%+v", v)))
			}
		default:
			// noop
		}
		return true
	})

	// If addCallerSkip is set, add caller skip to the record.
	if skip := s.addCallerSkip + contextAddCallerSkip(ctx); skip > 0 {
		var pcs [1]uintptr
		runtime.Callers(skip, pcs[:])
		r.PC = pcs[0]
	}

	// If addTimestamp is false, remove timestamp from the record.
	if !s.addTimestamp {
		r.Time = time.Time{}
	}

	// Add attrs to the record.
	if len(s.addAttrs) > 0 {
		r.AddAttrs(s.addAttrs...)
	}
	if len(attrs) > 0 {
		r.AddAttrs(attrs...)
	}

	//nolint:wrapcheck
	return s.slogHandler.Handle(ctx, r)
}

func (s *slogJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h := s.clone()
	h.slogHandler = h.slogHandler.WithAttrs(attrs)
	return h
}

func (s *slogJSONHandler) WithGroup(name string) slog.Handler {
	h := s.clone()
	h.slogHandler = h.slogHandler.WithGroup(name)
	return h
}

func (s *slogJSONHandler) clone() *slogJSONHandler {
	// Clone slogHandler
	c := *s
	// Clone slogHandlerOptions
	slogOptions := *s.slogHandlerOptions
	c.slogHandlerOptions = &slogOptions
	// Clone slogHandler
	c.slogHandler = s.slogHandler.WithAttrs(nil) // call WithAttrs with nil to clone slogHandler
	return &c
}
