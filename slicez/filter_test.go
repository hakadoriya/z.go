package slicez

import (
	"testing"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      []int
		shouldKeep func(int, int) bool
		expected   []int
	}{
		{
			name:  "success,keep_even_numbers",
			input: []int{1, 2, 3, 4, 5},
			shouldKeep: func(index int, elem int) bool {
				return elem%2 == 0
			},
			expected: []int{2, 4},
		},
		{
			name:  "success,keep_numbers_greater_than_3",
			input: []int{1, 2, 3, 4, 5},
			shouldKeep: func(index int, elem int) bool {
				return elem > 3
			},
			expected: []int{4, 5},
		},
		{
			name:  "success,keep_all_numbers",
			input: []int{1, 2, 3},
			shouldKeep: func(index int, elem int) bool {
				return true
			},
			expected: []int{1, 2, 3},
		},
		{
			name:  "success,keep_no_numbers",
			input: []int{1, 2, 3},
			shouldKeep: func(index int, elem int) bool {
				return false
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := Filter(tt.input, tt.shouldKeep)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected %d, got %d", tt.expected[i], v)
				}
			}
		})
	}
}
