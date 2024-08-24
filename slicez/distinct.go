package slicez

import "reflect"

func Distinct[T comparable](src []T) []T {
	for i := 0; i < len(src); i++ {
		for j := i + 1; j < len(src); j++ {
			if src[i] == src[j] {
				src = append(src[:j], src[j+1:]...)
				j--
			}
		}
	}

	return src
}

func DeepDistinct[T any](src []T) []T {
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
