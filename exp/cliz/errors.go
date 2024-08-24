package cliz

import "errors"

var (
	ErrDuplicateOption     = errors.New("duplicate option")
	ErrDuplicateSubCommand = errors.New("duplicate sub command")
	ErrHelp                = errors.New("help requested")
	ErrInvalidOptionType   = errors.New("invalid option type")
	ErrMissingOptionValue  = errors.New("missing option value")
	ErrNilContext          = errors.New("nil context")
	ErrNotSetInContext     = errors.New("not set in context")
	ErrNotCalled           = errors.New("not called")
	ErrOptionRequired      = errors.New("option required")
	ErrUnknownOption       = errors.New("unknown option")
)
