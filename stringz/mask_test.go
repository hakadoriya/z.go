package stringz

import (
	"testing"
)

func TestMaskPrefix(t *testing.T) {
	t.Parallel()
	type args struct {
		s         string
		mask      string
		unmaskLen int
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{"empty", args{"", "*", 0}, ""},
		{"mask_0", args{"ABCDEFGH", "*", 0}, "********"},
		{"mask_1", args{"ABCDEFGH", "*", 1}, "*******H"},
		{"mask_2", args{"ABCDEFGH", "*", 2}, "******GH"},
		{"mask_3", args{"ABCDEFGH", "*", 3}, "*****FGH"},
		{"mask_4", args{"ABCDEFGH", "*", 4}, "****EFGH"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if actual := MaskPrefix(tt.args.s, tt.args.mask, tt.args.unmaskLen); actual != tt.expected {
				t.Errorf("expected(%v) != actual(%v)", tt.expected, actual)
			}
		})
	}
}

func TestMaskSuffix(t *testing.T) {
	t.Parallel()
	type args struct {
		s         string
		mask      string
		unmaskLen int
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{"empty", args{"", "*", 0}, ""},
		{"mask_0", args{"ABCDEFGH", "*", 0}, "********"},
		{"mask_1", args{"ABCDEFGH", "*", 1}, "A*******"},
		{"mask_2", args{"ABCDEFGH", "*", 2}, "AB******"},
		{"mask_3", args{"ABCDEFGH", "*", 3}, "ABC*****"},
		{"mask_4", args{"ABCDEFGH", "*", 4}, "ABCD****"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if actual := MaskSuffix(tt.args.s, tt.args.mask, tt.args.unmaskLen); actual != tt.expected {
				t.Errorf("expected(%v) != actual(%v)", tt.expected, actual)
			}
		})
	}
}
