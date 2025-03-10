package slicez

func ForEach[T any](s []T, f func(index int, elem T)) {
	for idx, e := range s {
		f(idx, e)
	}
}
