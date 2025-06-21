package slogz

import "errors"

var ErrHandlerIsNotSlogJSONHandler = errors.New("slogz: slog.Handler is not *slogz.slogJSONHandler")
