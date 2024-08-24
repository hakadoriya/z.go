package otelz

import "reflect"

func deepDistinct[T any](src []T) []T {
	for i := 0; i < len(src); i++ {
		for j := i + 1; j < len(src); j++ {
			if reflect.DeepEqual(src[i], src[j]) {
				src = append(src[:j], src[j+1:]...)
				j--
			}
		}
	}

	return src
}

func reverse[T any](src []T) []T {
	rev := make([]T, 0, len(src))

	for i := range src {
		rev = append(rev, src[len(src)-1-i])
	}

	return rev
}

func isNil(v interface{}) bool {
	return (v == nil) || reflect.ValueOf(v).IsNil()
}
