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
	// ErrFieldCannotBeSet field cannot be set.
	ErrFieldCannotBeSet = errors.New("field cannot be set")
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
