package slicez

import (
	"strconv"
	"testing"
)

func TestFilterMap(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		generator func(index int, elem int) (string, bool)
		expected  []string
	}{
		{
			name:  "success,even_numbers",
			input: []int{1, 2, 3, 4, 5},
			generator: func(index int, elem int) (string, bool) {
				if elem%2 == 0 {
					return strconv.Itoa(elem), true
				}
				return "", false
			},
			expected: []string{"2", "4"},
		},
		{
			name:  "success,odd_numbers",
			input: []int{10, 15, 20, 25},
			generator: func(index int, elem int) (string, bool) {
				if elem%2 != 0 {
					return strconv.Itoa(elem), true
				}
				return "", false
			},
			expected: []string{"15", "25"},
		},
		{
			name:     "success,empty_slice",
			input:    []int{},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FilterMap(test.input, test.generator)
			if len(result) != len(test.expected) {
				t.Errorf("❌: expected length %d, got %d", len(test.expected), len(result))
			}

			for i, v := range result {
				if v != test.expected[i] {
					t.Errorf("❌: at index %d: expected %s, got %s", i, test.expected[i], v)
				}
			}
		})
	}
}
