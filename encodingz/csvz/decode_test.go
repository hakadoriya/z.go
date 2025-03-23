package csvz

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
	"time"
)

type testStruct struct {
	UserID     int    `csv:"user_id"`
	Username   string `csv:"username"`
	Hyphen     string `csv:"-"`
	Empty      string `csv:""`
	unexported string
	CreatedAt  time.Time `csv:"created_at"`
}

func TestCSVDecoder_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		argCSV       string
		expected     any
		actualTarget any
	}{
		{
			name:   "success,normal",
			argCSV: "user_id,username,created_at,etc\n0,user_0,2025-01-01T00:00:00Z,etc_0\n1,user_1,2025-01-02T00:00:00Z,etc_1",
			expected: &[]*testStruct{
				{UserID: 0, Username: "user_0", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
				{UserID: 1, Username: "user_1", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			actualTarget: &[]*testStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCSVDecoder(
				strings.NewReader(tt.argCSV),
				WithCSVDecoderOptionTagName("csv"),
				WithCSVDecoderOptionTimeFormat(time.RFC3339Nano),
				WithCSVDecoderOptionCSVReaderModifier(func(r *csv.Reader) *csv.Reader {
					r.Comma = ','
					return r
				}),
			)
			if err := c.Decode(tt.actualTarget); err != nil {
				t.Fatalf("❌: err != nil: %v", err)
			}
			if !reflect.DeepEqual(tt.actualTarget, tt.expected) {
				t.Errorf("❌: expected(%q) != actual(%q)", tt.expected, tt.actualTarget)
			}
		})
	}
}
