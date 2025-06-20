package csvz

import (
	"encoding/csv"
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
	csvReader         *csv.Reader
	tagName           string
	timeFormat        string
	setFieldValueFunc func(refrectType reflect.StructField, refrectValue reflect.Value, value string) (ok bool)

	headers  []string
	fieldMap map[string]int // Mapping from header names to indices
}

// NewCSVDecoder creates a new CSVDecoder
func NewCSVDecoder(r io.Reader, opts ...CSVDecoderOption) *CSVDecoder {
	d := &CSVDecoder{
		csvReader:         csv.NewReader(r),
		tagName:           defaultCSVTagName,
		timeFormat:        defaultTimeFormat,
		setFieldValueFunc: nil,

		headers:  make([]string, 0),
		fieldMap: make(map[string]int),
	}
	for _, opt := range opts {
		opt.apply(d)
	}
	return d
}

// Decode decodes CSV into Go structs
//
//nolint:cyclop
func (csvDecoder *CSVDecoder) Decode(v interface{}) error {
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
		return ErrDecodeTargetMustBeStructSlice
	}

	// Read headers
	headers, err := csvDecoder.csvReader.Read()
	if err != nil {
		return fmt.Errorf("csvDecoder.csvReader.Read: headers=%v: %w", headers, err)
	}
	csvDecoder.headers = headers

	// Build mapping between headers and fields
	for i, header := range headers {
		csvDecoder.fieldMap[header] = i
	}

	// Read each row and convert to struct
	for {
		record, err := csvDecoder.csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("csvDecoder.csvReader.Read: record=%v: %w", record, err)
		}

		// Create new struct instance
		newElem := reflect.New(elemType).Elem()

		// Map fields and set data
		if err := csvDecoder.mapFields(newElem, record); err != nil {
			return fmt.Errorf("csvDecoder.mapFields: %w", err)
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
func (csvDecoder *CSVDecoder) mapFields(structVal reflect.Value, record []string) error {
	structType := structVal.Type()

	// Process each field in the struct
	for i := range structVal.NumField() {
		fieldValue := structVal.Field(i)
		fieldType := structType.Field(i)

		// skip if field is private
		if !fieldType.IsExported() {
			continue
		}

		// Get column name from tag
		tag := fieldType.Tag.Get(csvDecoder.tagName)
		if tag == "" || tag == "-" {
			continue
		}

		// Get CSV column index corresponding to the tag
		colIdx, ok := csvDecoder.fieldMap[tag]
		if !ok {
			continue // Skip if no matching column
		}

		// Verify value is within range
		if colIdx >= len(record) {
			continue
		}

		// Convert string to appropriate type and set field
		if err := csvDecoder.setFieldValue(fieldType, fieldValue, record[colIdx]); err != nil {
			return fmt.Errorf("csvDecoder.setFieldValue: name=%s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue converts a string value to the appropriate type and sets it to the field
//
//nolint:cyclop,funlen,gocognit
func (csvDecoder *CSVDecoder) setFieldValue(refrectType reflect.StructField, refrectValue reflect.Value, value string) error {
	if csvDecoder.setFieldValueFunc != nil {
		if ok := csvDecoder.setFieldValueFunc(refrectType, refrectValue, value); ok {
			return nil
		}
	}

	if !refrectValue.CanSet() {
		return ErrStructFieldCannotBeSet
	}

	if value == "" {
		// Keep default value for empty values
		return nil
	}

	// Check for CSVUnmarshaler interface implementation
	//
	//nolint:nestif
	if refrectValue.CanAddr() {
		pv := refrectValue.Addr()
		if pv.CanInterface() {
			if u, ok := pv.Interface().(CSVUnmarshaler); ok {
				if err := u.UnmarshalCSV(value); err != nil {
					return fmt.Errorf("error unmarshalling field %s: %w", refrectValue.Type().Name(), err)
				}
				return nil
			}
		}
	}

	const bitSize = 64
	//nolint:exhaustive // for testing
	switch refrectValue.Kind() {
	case reflect.String:
		refrectValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, bitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt: %w", err)
		}
		refrectValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, bitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint: %w", err)
		}
		refrectValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseFloat: %w", err)
		}
		refrectValue.SetFloat(f)
	case reflect.Complex64, reflect.Complex128:
		c, err := strconv.ParseComplex(value, bitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseComplex: %w", err)
		}
		refrectValue.SetComplex(c)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("strconv.ParseBool: %w", err)
		}
		refrectValue.SetBool(b)
	case reflect.Struct:
		// time.Time's special handling
		if refrectValue.Type() == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(csvDecoder.timeFormat, value)
			if err != nil {
				return fmt.Errorf("time.Parse: %w", err)
			}
			refrectValue.Set(reflect.ValueOf(t))
		} else {
			return fmt.Errorf("type=%s: %w", refrectValue.Type().Name(), ErrUnsupportedType)
		}
	case reflect.Pointer:
		if refrectValue.IsNil() {
			refrectValue.Set(reflect.New(refrectValue.Type().Elem()))
		}
		if err := csvDecoder.setFieldValue(refrectType, refrectValue.Elem(), value); err != nil {
			return fmt.Errorf("csvDecoder.setFieldValue: name=%s: %w", refrectType.Name, err)
		}
	// NOTE: for testing
	// case reflect.Invalid,
	// 	reflect.Uintptr,
	// 	reflect.Array,
	// 	reflect.Chan,
	// 	reflect.Func,
	// 	reflect.Interface,
	// 	reflect.Map,
	// 	reflect.Slice,
	// 	reflect.UnsafePointer:
	// 	return fmt.Errorf("kind=%s: %w", field.Kind(), ErrUnsupportedType)
	default:
		return fmt.Errorf("kind=%s: %w", refrectValue.Kind(), ErrUnsupportedType)
	}

	return nil
}
