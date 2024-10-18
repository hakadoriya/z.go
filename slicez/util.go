package slicez

func First[T any](src []T) T { //nolint:ireturn
	if len(src) == 0 {
		var zero T
		return zero
	}
	return src[0]
}

func Last[T any](src []T) T { //nolint:ireturn
	if len(src) == 0 {
		var zero T
		return zero
	}
	return src[len(src)-1]
}
