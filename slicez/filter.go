package slicez

func Filter[T any](s []T, shouldKeep func(index int, elem T) bool) []T {
	gen := make([]T, 0, len(s))
	for idx, e := range s {
		if shouldKeep(idx, e) {
			gen = append(gen, e)
		}
	}
	return gen
}
