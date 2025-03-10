package slicez

func Map[T, U any](src []T, generator func(index int, elem T) U) []U {
	b := make([]U, len(src))
	for i := range src {
		b[i] = generator(i, src[i])
	}
	return b
}
