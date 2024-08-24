package ctxerrz

import "errors"

var (
	ErrNilContext        = errors.New("nil context")
	ErrNotFoundInContext = errors.New("not found in context")
)
