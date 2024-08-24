package slicez

import "slices"

// CompactStable returns a new slice with all duplicate elements removed.
// The order of the elements is preserved by first occurrence.
func CompactStable[T comparable](a []T) []T {
	_uniq := make(map[T]struct{})
	const reasonableCap = 8
	uniq := make([]T, 0, reasonableCap)
	for i := range a {
		if _, ok := _uniq[a[i]]; !ok {
			uniq = append(uniq, a[i])
			_uniq[a[i]] = struct{}{}
		}
	}

	return slices.Clip(uniq)
}
