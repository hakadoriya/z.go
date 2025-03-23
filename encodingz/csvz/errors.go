package csvz

import "errors"

// decode
var (
	// ErrDecodeTargetMustBeNonNilPointer decode target must be a non-nil pointer.
	ErrDecodeTargetMustBeNonNilPointer = errors.New("decode target must be a non-nil pointer")
	// ErrDecodeTargetMustBeSlice decode target must be a slice.
	ErrDecodeTargetMustBeSlice = errors.New("decode target must be a slice")
	// ErrDecodeTargetMustBeStruct decode target must be a struct.
	ErrDecodeTargetMustBeStruct = errors.New("decode target must be a struct")
	// ErrStructFieldCannotBeSet struct field cannot be set.
	ErrStructFieldCannotBeSet = errors.New("struct field cannot be set; unexported field or field is not settable")
	// ErrUnsupportedType unsupported type.
	ErrUnsupportedType = errors.New("unsupported type")
)

// encode
var (
	// ErrEncodeSourceMustBeSlice encode source must be a slice.
	ErrEncodeSourceMustBeSlice = errors.New("encode source must be a slice")
	// ErrEncodeSourceMustBeStruct slice elements must be structs or pointers to structs.
	ErrEncodeSourceMustBeStruct = errors.New("slice elements must be structs or pointers to structs")
)
