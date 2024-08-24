package slicez

import "reflect"

func Contains[T comparable](s []T, value T) bool {
	for _, elem := range s {
		if value == elem {
			return true
		}
	}

	return false
}

func DeepContains[T any](s []T, value T) bool {
	for _, elem := range s {
		if reflect.DeepEqual(value, elem) {
			return true
		}
	}

	return false
}
