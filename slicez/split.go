package slicez

func Split[T any](src []T, size int) [][]T {
	sourceLength := len(src)
	splits := make([][]T, 0, sourceLength/size+1)
	for i := 0; i < sourceLength; i += size {
		next := i + size
		if sourceLength < next {
			next = sourceLength
		}
		splits = append(splits, src[i:next])
	}
	return splits
}
