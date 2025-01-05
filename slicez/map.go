package slicez

func Map[T, U any](src []T, f func(T) U) []U {
	b := make([]U, len(src))
	for i := range src {
		b[i] = f(src[i])
	}
	return b
}
