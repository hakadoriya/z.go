package csvz

import (
	"reflect"
	"strings"
	"testing"
)

func TestCSVDecoder_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		csv       string
		gotTarget interface{}
		expected  interface{}
	}{
		{name: "success,int_slice", csv: "1,2,3", gotTarget: &[]int{}, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCSVDecoder(strings.NewReader(tt.csv))
			if err := c.Decode(tt.gotTarget); err != nil {
				t.Fatalf("❌: err != nil: %v", err)
			}
			if !reflect.DeepEqual(tt.gotTarget, tt.expected) {
				t.Errorf("❌: expected(%q) != actual(%q)", tt.expected, tt.gotTarget)
			}
		})
	}
}
