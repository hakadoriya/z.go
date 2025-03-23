package csvz

import (
	"bytes"
	"encoding/csv"
	"testing"
	"time"
)

func TestCSVEncoder_Encode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		argSlice any
		expected string
	}{
		{
			name: "success,normal",
			argSlice: &[]*testStruct{
				{UserID: 0, Username: "user_0", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
				{UserID: 1, Username: "user_1", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			expected: "user_id,username,created_at\n0,user_0,2025-01-01T00:00:00Z\n1,user_1,2025-01-02T00:00:00Z\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			c := NewCSVEncoder(
				w,
				WithCSVEncoderOptionTagName("csv"),
				WithCSVEncoderOptionTimeFormat(time.RFC3339Nano),
				WithCSVEncoderOptionCSVWriterModifier(func(w *csv.Writer) *csv.Writer {
					w.Comma = ','
					return w
				}),
			)
			if err := c.Encode(tt.argSlice); err != nil {
				t.Fatalf("❌: err != nil: %v", err)
			}
			if expected, actual := tt.expected, w.String(); expected != actual {
				t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
			}
		})
	}
}
