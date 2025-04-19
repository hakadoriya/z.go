package csvz

import (
	"encoding/csv"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

type testStruct struct {
	UserID     uint    `csv:"user_id"`
	Username   string  `csv:"username"`
	Age        int     `csv:"age"`
	IsActive   bool    `csv:"is_active"`
	Point      float64 `csv:"point"`
	Hyphen     string  `csv:"-"`
	Empty      string  `csv:""`
	unexported string
	CreatedAt  time.Time `csv:"created_at"`
}

func TestCSVDecoder_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		argCSV      string
		requireFunc func(t *testing.T, name string, err error)
		argTarget   any
		expected    any
	}{
		{
			name:   "success,normal",
			argCSV: "user_id,username,age,is_active,point,created_at,etc\n0,\"user,0\",20,true,100.1,2025-01-01T00:00:00Z,etc_0\n1,user_1,21,false,200.2,2025-01-02T00:00:00Z,etc_1",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			argTarget: &[]*testStruct{},
			expected: &[]*testStruct{
				{UserID: 0, Username: "user,0", Age: 20, IsActive: true, Point: 100.1, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
				{UserID: 1, Username: "user_1", Age: 21, IsActive: false, Point: 200.2, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:   "failure,ErrDecodeTargetMustBeNonNilPointer",
			argCSV: "user_id,username,age,is_active,point,created_at,etc\n0,\"user,0\",20,true,100.1,2025-01-01T00:00:00Z,etc_0\n1,user_1,21,false,200.2,2025-01-02T00:00:00Z,etc_1",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrDecodeTargetMustBeNonNilPointer) {
					t.Fatalf("❌: name=%s: err != ErrDecodeTargetMustBeNonNilPointer: %v", name, err)
				}
			},
			argTarget: nil,
			expected:  nil,
		},
		{
			name:   "failure,ErrDecodeTargetMustBeSlice",
			argCSV: "user_id,username,age,is_active,point,created_at,etc\n0,\"user,0\",20,true,100.1,2025-01-01T00:00:00Z,etc_0\n1,user_1,21,false,200.2,2025-01-02T00:00:00Z,etc_1",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrDecodeTargetMustBeSlice) {
					t.Fatalf("❌: name=%s: err != ErrDecodeTargetMustBeSlice: %v", name, err)
				}
			},
			argTarget: &struct{}{},
			expected:  &struct{}{},
		},
		{
			name:   "failure,ErrDecodeTargetMustBeStructSlice",
			argCSV: "user_id,username,age,is_active,point,created_at,etc\n0,\"user,0\",20,true,100.1,2025-01-01T00:00:00Z,etc_0\n1,user_1,21,false,200.2,2025-01-02T00:00:00Z,etc_1",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrDecodeTargetMustBeStructSlice) {
					t.Fatalf("❌: name=%s: err != ErrDecodeTargetMustBeStructSlice: %v", name, err)
				}
			},
			argTarget: &[]int{},
			expected:  &[]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewCSVDecoder(
				strings.NewReader(tt.argCSV),
				WithCSVDecoderOptionTagName("csv"),
				WithCSVDecoderOptionTimeFormat(time.RFC3339Nano),
				WithCSVDecoderOptionCSVReaderModifier(func(r *csv.Reader) *csv.Reader {
					r.Comma = ','
					return r
				}),
			)
			err := c.Decode(tt.argTarget)
			tt.requireFunc(t, tt.name, err)
			if !reflect.DeepEqual(tt.expected, tt.argTarget) {
				t.Errorf("❌: name=%s: expected(%q) != actual(%q)", tt.name, tt.expected, tt.argTarget)
			}
		})
	}
}

func TestCSVDecoder_setFieldValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		argCSVDecoder *CSVDecoder
		argField      reflect.Value
		argValue      string
		requireFunc   func(t *testing.T, name string, err error)
		expected      any
	}{
		{
			name:          "failure,ErrStructFieldCannotBeSet",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(testStruct{})),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrStructFieldCannotBeSet) {
					t.Fatalf("❌: name=%s: err != ErrStructFieldCannotBeSet: %v", name, err)
				}
			},
			expected: &testStruct{},
		},
		{
			name:          "failure,int64",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(int64(0))).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "strconv.ParseInt") {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: int64(0),
		},
		{
			name:          "failure,uint64",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(uint64(0))).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "strconv.ParseUint") {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: uint64(0),
		},
		{
			name:          "failure,float64",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(float64(0))).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "strconv.ParseFloat") {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: float64(0),
		},
		{
			name:          "failure,bool",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(false)).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "strconv.ParseBool") {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: false,
		},
		{
			name:          "failure,time.Time",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(time.Time{})).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "time.Parse") {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: time.Time{},
		},
		{
			name:          "success,complex128",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(complex(0, 0))).Elem(),
			argValue:      "100i",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: complex(0, 100),
		},
		{
			name:          "failure,complex128",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(complex(0, 0))).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "strconv.ParseComplex") {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: complex(0, 0),
		},
		{
			name:          "failure,ErrUnsupportedType",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(struct{}{})).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrUnsupportedType) {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: struct{}{},
		},
		{
			name:          "failure,ErrUnsupportedType",
			argCSVDecoder: &CSVDecoder{},
			argField:      reflect.New(reflect.TypeOf(chan int(nil))).Elem(),
			argValue:      "a",
			requireFunc: func(t *testing.T, name string, err error) {
				t.Helper()
				if !errors.Is(err, ErrUnsupportedType) {
					t.Fatalf("❌: name=%s: err != nil: %v", name, err)
				}
			},
			expected: chan int(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.argCSVDecoder.setFieldValue(tt.argField, tt.argValue)
			tt.requireFunc(t, tt.name, err)
			if !reflect.DeepEqual(tt.expected, tt.argField.Interface()) {
				t.Errorf("❌: name=%s: expected(%q) != actual(%q)", tt.name, tt.expected, tt.argField.Interface())
			}
		})
	}
}

// TODO: Add test
