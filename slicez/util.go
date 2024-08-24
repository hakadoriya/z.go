package slicez

func First[T any](s []T) T { //nolint:ireturn
	if len(s) == 0 {
		var zero T
		return zero
	}
	return s[0]
}

func Last[T any](s []T) T { //nolint:ireturn
	if len(s) == 0 {
		var zero T
		return zero
	}
	return s[len(s)-1]
}
