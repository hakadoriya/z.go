package slicez

import "reflect"

func Contains[T comparable](src []T, value T) bool {
	for _, elem := range src {
		if value == elem {
			return true
		}
	}

	return false
}

func DeepContains[T any](src []T, value T) bool {
	for _, elem := range src {
		if reflect.DeepEqual(value, elem) {
			return true
		}
	}

	return false
}
