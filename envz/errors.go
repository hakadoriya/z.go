package envz

import (
	"errors"
	"strconv"
)

var (
	ErrRange                                     = errors.New(strconv.ErrRange.Error())
	ErrEnvironmentVariableIsEmpty                = errors.New("environment variable is empty")
	ErrInvalidType                               = errors.New("invalid type; must be a pointer to a struct")
	ErrStructFieldCannotBeSet                    = errors.New("struct field cannot be set; unexported field or field is not settable")
	ErrInvalidTagValueEnvironmentVariableIsEmpty = errors.New("invalid tag value; environment variable name is empty")
	ErrInvalidTagValueInvalidKey                 = errors.New("invalid tag value; invalid key")
	ErrRequiredEnvironmentVariableNotFound       = errors.New("required environment variable not found")
	ErrStructFieldTypeNotSupported               = errors.New("struct field type not supported")
)
