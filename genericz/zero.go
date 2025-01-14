package genericz

func Zero[T interface{}](_ T) (zero T) { return zero } //nolint:ireturn

func SliceElemZero[T interface{}](_ []T) (zero T) { return zero } //nolint:ireturn

func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}
