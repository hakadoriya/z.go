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
	csvWriter  *csv.Writer
	tagName    string
	timeFormat string
}

// NewCSVEncoder creates a new CSVEncoder
func NewCSVEncoder(w io.Writer, opts ...CSVEncoderOption) *CSVEncoder {
	e := &CSVEncoder{
		csvWriter:  csv.NewWriter(w),
		tagName:    defaultCSVTagName,
		timeFormat: defaultTimeFormat,
	}
	for _, opt := range opts {
		opt.apply(e)
	}
	return e
}

// Encode encodes Go structs to CSV
//
//nolint:cyclop
func (csvEncoder *CSVEncoder) Encode(v interface{}) error {
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
		return ErrEncodeSourceMustBeStructSlice
	}

	// Generate headers
	headers, fieldIndices := csvEncoder.extractHeaders(elemVal.Type())
	if err := csvEncoder.csvWriter.Write(headers); err != nil {
		return fmt.Errorf("headers: csvEncoder.csvWriter.Write: %w", err)
	}

	// Encode each row
	for i := range rv.Len() {
		rowVal := rv.Index(i)
		if rowVal.Kind() == reflect.Ptr {
			rowVal = rowVal.Elem()
		}

		record := make([]string, len(headers))
		for j, idx := range fieldIndices {
			if idx >= 0 {
				field := rowVal.Field(idx)
				record[j] = csvEncoder.fieldToString(field)
			}
		}

		if err := csvEncoder.csvWriter.Write(record); err != nil {
			return fmt.Errorf("row=%d: csvEncoder.csvWriter.Write: %w", i, err)
		}
	}

	csvEncoder.csvWriter.Flush()
	if err := csvEncoder.csvWriter.Error(); err != nil {
		return fmt.Errorf("csvEncoder.csvWriter.Error: %w", err)
	}
	return nil
}

// extractHeaders extracts CSV headers from struct tags
func (csvEncoder *CSVEncoder) extractHeaders(t reflect.Type) ([]string, []int) {
	headers := make([]string, 0, t.NumField())
	fieldMap := make(map[string]int)

	// First collect all fields with tags
	for i := range t.NumField() {
		fieldType := t.Field(i)

		// skip if field is private
		if !fieldType.IsExported() {
			continue
		}

		// Get column name from tag
		tag := fieldType.Tag.Get(csvEncoder.tagName)
		if tag == "" || tag == "-" {
			continue
		}

		fieldMap[tag] = i
		headers = append(headers, tag)
	}

	// Prepare field indices corresponding to headers
	fieldIndices := make([]int, len(headers))
	for i, header := range headers {
		fieldIndices[i] = fieldMap[header]
	}

	return headers, fieldIndices
}

// fieldToString converts a field value to string
//
//nolint:cyclop
func (csvEncoder *CSVEncoder) fieldToString(field reflect.Value) string {
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

	//nolint:exhaustive // for testing
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
			//nolint:forcetypeassert
			return field.Interface().(time.Time).Format(csvEncoder.timeFormat)
		}
		return fmt.Sprintf("%v", field.Interface())
	// NOTE: for testing
	// case reflect.Invalid,
	// 	reflect.Uintptr,
	// 	reflect.Complex64,
	// 	reflect.Complex128,
	// 	reflect.Array,
	// 	reflect.Chan,
	// 	reflect.Func,
	// 	reflect.Interface,
	// 	reflect.Map,
	// 	reflect.Pointer,
	// 	reflect.Slice,
	// 	reflect.UnsafePointer:
	// 	return fmt.Sprintf("%v", field.Interface())
	default:
		return fmt.Sprintf("%v", field.Interface())
	}
}
