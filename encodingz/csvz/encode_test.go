package csvz

import (
	"bytes"
	"encoding/csv"
	"errors"
	"testing"
	"time"
)

func TestCSVEncoder_Encode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        []CSVEncoderOption
		argSlice    any
		requireFunc func(t *testing.T, name string, err error)
		expected    string
	}{
		{
			name: "success,normal",
			argSlice: &[]*testStruct{
				{UserID: 0, Username: "user,0", Age: 20, IsActive: true, Point: 100.1, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
				{UserID: 1, Username: "user_1", Age: 21, IsActive: false, Point: 200.2, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: "user_id,username,age,is_active,point,created_at\n0,\"user,0\",20,true,100.1,2025-01-01T00:00:00Z\n1,user_1,21,false,200.2,2025-01-02T00:00:00Z\n",
		},
		{
			name:     "success,empty_slice",
			argSlice: &[]*testStruct{},
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: "",
		},
		{
			name:     "failure,ErrEncodeSourceMustBeSlice",
			argSlice: &struct{}{},
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrEncodeSourceMustBeSlice) {
					t.Fatalf("❌: name=%s: err != ErrEncodeSourceMustBeSlice: %v", name, err)
				}
			},
		},
		{
			name:     "failure,ErrEncodeSourceMustBeStructSlice",
			argSlice: &[]int{0, 1},
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrEncodeSourceMustBeStructSlice) {
					t.Fatalf("❌: name=%s: err != ErrEncodeSourceMustBeStructSlice: %v", name, err)
				}
			},
		},
		{
			name: "failure,csvEncoder.csvWriter.Write",
			opts: []CSVEncoderOption{
				WithCSVEncoderOptionCSVWriterModifier(func(w *csv.Writer) *csv.Writer {
					w.Comma = rune(0)
					return w
				}),
			},
			argSlice: &[]*testStruct{
				{UserID: 0, Username: "user,0", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
				{UserID: 1, Username: "user_1", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err == nil {
					t.Fatalf("❌: name=%s: err == nil", name)
				}
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := bytes.NewBuffer(nil)
			c := NewCSVEncoder(
				w,
				append(
					[]CSVEncoderOption{
						WithCSVEncoderOptionTagName("csv"),
						WithCSVEncoderOptionTimeFormat(time.RFC3339Nano),
						WithCSVEncoderOptionCSVWriterModifier(func(w *csv.Writer) *csv.Writer {
							w.Comma = ','
							return w
						}),
					},
					tt.opts...,
				)...,
			)
			err := c.Encode(tt.argSlice)
			tt.requireFunc(t, tt.name, err)
			if expected, actual := tt.expected, w.String(); expected != actual {
				t.Errorf("❌: name=%s: expected(%q) != actual(%q)", tt.name, expected, actual)
			}
		})
	}
}

// TODO: Add test
