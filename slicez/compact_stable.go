package slicez

import "slices"

// CompactStable returns a new slice with all duplicate elements removed.
// The order of the elements is preserved by first occurrence.
func CompactStable[T comparable](src []T) []T {
	_uniq := make(map[T]struct{})

	const defaultCap = 8
	uniq := make([]T, 0, defaultCap)
	for i := range src {
		if _, ok := _uniq[src[i]]; !ok {
			uniq = append(uniq, src[i])
			_uniq[src[i]] = struct{}{}
		}
	}

	return slices.Clip(uniq)
}
