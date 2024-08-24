package clicorez

import (
	"errors"

	"github.com/hakadoriya/z.go/contextz/ctxerrz"
)

var (
	ErrDuplicateOption             = errors.New("duplicate option")
	ErrDuplicateSubCommand         = errors.New("duplicate sub command")
	ErrHelp                        = errors.New("help")
	ErrInvalidOptionType           = errors.New("invalid option type")
	ErrMissingOptionValue          = errors.New("missing option value")
	ErrNilContext                  = ctxerrz.ErrNilContext
	ErrNotFoundInContext           = ctxerrz.ErrNotFoundInContext
	ErrNotCalled                   = errors.New("not called")
	ErrOptionRequired              = errors.New("option required")
	ErrUnknownOption               = errors.New("unknown option")
	ErrInvalidType                 = errors.New("invalid type; must be a pointer to a struct")
	ErrStructFieldCannotBeSet      = errors.New("struct field cannot be set; unexported field or field is not settable")
	ErrInvalidTagValue             = errors.New("invalid tag value")
	ErrStructFieldTypeNotSupported = errors.New("struct field type not supported")
)
