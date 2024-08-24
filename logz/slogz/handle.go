package slogz

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
)

var _ slog.Handler = (*slogHandler)(nil)

type slogHandler struct {
	slogHandlerOptions *slog.HandlerOptions

	errorVerbose          bool
	errorVerboseKeySuffix string
	slogHandler           slog.Handler
}

func NewHandler(w io.Writer, level slog.Leveler, opts ...HandlerOption) slog.Handler {
	s := &slogHandler{
		slogHandlerOptions: &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: ReplaceAttr,
		},
		errorVerbose:          true,
		errorVerboseKeySuffix: "Verbose",
	}

	s.slogHandler = slog.NewJSONHandler(w, s.slogHandlerOptions)

	return s
}

func (s *slogHandler) clone() *slogHandler {
	// Clone slogHandler
	c := *s
	// Clone slogHandlerOptions
	slogOptions := *s.slogHandlerOptions
	c.slogHandlerOptions = &slogOptions
	// Clone slogHandler
	// XXX: There is no way to clone?
	c.slogHandler = s.slogHandler.WithAttrs(nil)
	return &c
}

func (s *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.slogHandler.Enabled(ctx, level)
}

func (s *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		// Attr Value type switch
		switch v := a.Value.Any().(type) {
		case error:
			// If errorVerbose is set, add verbose error to the record.
			if s.errorVerbose {
				attrs = append(attrs, slog.String(a.Key+s.errorVerboseKeySuffix, fmt.Sprintf("%+v", v)))
			}
		}
		return true
	})

	// If AddCallerSkip is set, add caller skip to the record.
	if skip := addCallerSkip(ctx); skip > 0 {
		const defaultCallerSkip = 4
		var pcs [1]uintptr
		runtime.Callers(defaultCallerSkip+skip, pcs[:])
		r.PC = pcs[0]
	}

	// Add attrs to the record.
	if len(attrs) > 0 {
		r.AddAttrs(attrs...)
	}

	return s.slogHandler.Handle(ctx, r)
}

func (s *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return s.clone().slogHandler.WithAttrs(attrs)
}

func (s *slogHandler) WithGroup(name string) slog.Handler {
	return s.clone().slogHandler.WithGroup(name)
}
