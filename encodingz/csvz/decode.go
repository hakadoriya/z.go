package csvz

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

// CSVUnmarshaler is an interface that defines how custom types unmarshal from CSV
type CSVUnmarshaler interface {
	UnmarshalCSV(value string) error
}

// CSVDecoder is a decoder for decoding CSV into Go structs
type CSVDecoder struct {
	r          *csv.Reader
	tagName    string
	timeFormat string
	headers    []string
	fieldMap   map[string]int // Mapping from header names to indices
}

// NewCSVDecoder creates a new CSVDecoder
func NewCSVDecoder(r io.Reader, opts ...CSVDecoderOption) *CSVDecoder {
	d := &CSVDecoder{
		r:        csv.NewReader(r),
		tagName:  defaultCSVTagName,
		fieldMap: make(map[string]int),
	}
	for _, opt := range opts {
		opt.apply(d)
	}
	return d
}

// Decode decodes CSV into Go structs
func (d *CSVDecoder) Decode(v interface{}) error {
	// Verify it's a pointer
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrDecodeTargetMustBeNonNilPointer
	}

	// Verify it's a slice
	sliceVal := rv.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return ErrDecodeTargetMustBeSlice
	}

	// Get the slice element type
	elemType := sliceVal.Type().Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return ErrDecodeTargetMustBeStruct
	}

	// Read headers
	headers, err := d.r.Read()
	if err != nil {
		return err
	}
	d.headers = headers

	// Build mapping between headers and fields
	for i, header := range headers {
		d.fieldMap[header] = i
	}

	// Read each row and convert to struct
	for {
		record, err := d.r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Create new struct instance
		newElem := reflect.New(elemType).Elem()

		// Map fields and set data
		if err := d.mapFields(newElem, record); err != nil {
			return err
		}

		// Append to slice
		if isPtr {
			sliceVal.Set(reflect.Append(sliceVal, newElem.Addr()))
		} else {
			sliceVal.Set(reflect.Append(sliceVal, newElem))
		}
	}

	return nil
}

// mapFields maps CSV values to struct fields
func (d *CSVDecoder) mapFields(structVal reflect.Value, record []string) error {
	structType := structVal.Type()

	// Process each field in the struct
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)

		// Get column name from tag
		tag := fieldType.Tag.Get(d.tagName)
		if tag == "" || tag == "-" {
			continue
		}

		// Get CSV column index corresponding to the tag
		colIdx, ok := d.fieldMap[tag]
		if !ok {
			continue // Skip if no matching column
		}

		// Verify value is within range
		if colIdx >= len(record) {
			continue
		}

		// Convert string to appropriate type and set field
		if err := d.setFieldValue(field, record[colIdx]); err != nil {
			return fmt.Errorf("error setting field %s: %v", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue converts a string value to the appropriate type and sets it to the field
func (d *CSVDecoder) setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return errors.New("field cannot be set")
	}

	if value == "" {
		// Keep default value for empty values
		return nil
	}

	// Check for CSVUnmarshaler interface implementation
	if field.CanAddr() {
		pv := field.Addr()
		if pv.CanInterface() {
			if u, ok := pv.Interface().(CSVUnmarshaler); ok {
				return u.UnmarshalCSV(value)
			}
		}
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Struct:
		// time.Time's special handling
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(d.timeFormat, value)
			if err != nil {
				return fmt.Errorf("time.Parse: %w", err)
			}
			field.Set(reflect.ValueOf(t))
		} else {
			return fmt.Errorf("type=%s: %w", field.Type().Name(), ErrUnsupportedType)
		}
	default:
		return fmt.Errorf("kind=%s: %w", field.Kind(), ErrUnsupportedType)
	}

	return nil
}
