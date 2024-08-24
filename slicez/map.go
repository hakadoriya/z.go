package slicez

func Map[T, U any](a []T, f func(T) U) []U {
	b := make([]U, len(a))
	for i := range a {
		b[i] = f(a[i])
	}
	return b
}
