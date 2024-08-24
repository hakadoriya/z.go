package envcorez

import "errors"

var (
	ErrInvalidType                         = errors.New("invalid type; must be a pointer to a struct")
	ErrStructFieldCannotBeSet              = errors.New("struct field cannot be set; unexported field or field is not settable")
	ErrInvalidTagValue                     = errors.New("invalid tag value")
	ErrRequiredEnvironmentVariableNotFound = errors.New("required environment variable not found")
	ErrStructFieldTypeNotSupported         = errors.New("struct field type not supported")
)
