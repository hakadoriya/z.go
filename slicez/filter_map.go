package slicez

func FilterMap[T, U any](s []T, generator func(index int, elem T) (v U, shouldKeep bool)) []U {
	gen := make([]U, 0, len(s))
	for idx, e := range s {
		if v, shouldKeep := generator(idx, e); shouldKeep {
			gen = append(gen, v)
		}
	}
	return gen
}
