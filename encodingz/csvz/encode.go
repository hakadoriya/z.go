package csvz

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

// CSVMarshaler is an interface that defines how custom types marshal to CSV
type CSVMarshaler interface {
	MarshalCSV() (value string, err error)
}

// CSVEncoder is an encoder for encoding Go structs to CSV
type CSVEncoder struct {
	w          *csv.Writer
	tagName    string
	timeFormat string
}

// NewCSVEncoder creates a new CSVEncoder
func NewCSVEncoder(w io.Writer, opts ...CSVEncoderOption) *CSVEncoder {
	e := &CSVEncoder{
		w:       csv.NewWriter(w),
		tagName: defaultCSVTagName,
	}
	for _, opt := range opts {
		opt.apply(e)
	}
	return e
}

// Encode encodes Go structs to CSV
func (e *CSVEncoder) Encode(v interface{}) error {
	rv := reflect.ValueOf(v)

	// Handle single object if not a slice
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Slice {
		return ErrEncodeSourceMustBeSlice
	}

	if rv.Len() == 0 {
		return nil // Do nothing for empty slice
	}

	// Get struct type information from first element
	elemVal := rv.Index(0)
	if elemVal.Kind() == reflect.Ptr {
		elemVal = elemVal.Elem()
	}

	if elemVal.Kind() != reflect.Struct {
		return ErrEncodeSourceMustBeStruct
	}

	// Generate headers
	headers, fieldIndices := e.extractHeaders(elemVal.Type())
	if err := e.w.Write(headers); err != nil {
		return err
	}

	// Encode each row
	for i := 0; i < rv.Len(); i++ {
		rowVal := rv.Index(i)
		if rowVal.Kind() == reflect.Ptr {
			rowVal = rowVal.Elem()
		}

		record := make([]string, len(headers))
		for j, idx := range fieldIndices {
			if idx >= 0 {
				field := rowVal.Field(idx)
				record[j] = e.fieldToString(field)
			}
		}

		if err := e.w.Write(record); err != nil {
			return err
		}
	}

	e.w.Flush()
	return e.w.Error()
}

// extractHeaders extracts CSV headers from struct tags
func (e *CSVEncoder) extractHeaders(t reflect.Type) ([]string, []int) {
	var headers []string
	var fieldIndices []int

	fieldMap := make(map[string]int)

	// First collect all fields with tags
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(e.tagName)
		if tag == "" || tag == "-" {
			continue
		}

		fieldMap[tag] = i
		headers = append(headers, tag)
	}

	// Prepare field indices corresponding to headers
	fieldIndices = make([]int, len(headers))
	for i, header := range headers {
		fieldIndices[i] = fieldMap[header]
	}

	return headers, fieldIndices
}

// fieldToString converts a field value to string
func (e *CSVEncoder) fieldToString(field reflect.Value) string {
	// Check for CSVMarshaler interface implementation
	if field.CanInterface() {
		if m, ok := field.Interface().(CSVMarshaler); ok {
			str, err := m.MarshalCSV()
			if err == nil {
				return str
			}
			// Fallback if error occurs
		}
	}

	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(field.Bool())
	case reflect.Struct:
		// Special handling for time.Time
		if field.Type() == reflect.TypeOf(time.Time{}) {
			return field.Interface().(time.Time).Format(e.timeFormat)
		}
		return fmt.Sprintf("%v", field.Interface())
	default:
		return fmt.Sprintf("%v", field.Interface())
	}
}
