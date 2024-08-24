package slicez

func Reverse[T any](src []T) []T {
	rev := make([]T, 0, len(src))

	for i := range src {
		rev = append(rev, src[len(src)-1-i])
	}

	return rev
}
