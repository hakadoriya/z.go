package csvz

import (
	"encoding/csv"
	"reflect"
)

type (
	CSVDecoderOption     interface{ apply(d *CSVDecoder) }
	csvDecoderOptionFunc func(d *CSVDecoder)
)

func (f csvDecoderOptionFunc) apply(d *CSVDecoder) { f(d) }

func WithCSVDecoderOptionCSVReaderModifier(modifier func(r *csv.Reader) *csv.Reader) CSVDecoderOption {
	return csvDecoderOptionFunc(func(d *CSVDecoder) { d.csvReader = modifier(d.csvReader) })
}

func WithCSVDecoderOptionTagName(tagName string) CSVDecoderOption {
	return csvDecoderOptionFunc(func(d *CSVDecoder) { d.tagName = tagName })
}

func WithCSVDecoderOptionTimeFormat(timeFormat string) CSVDecoderOption {
	return csvDecoderOptionFunc(func(d *CSVDecoder) { d.timeFormat = timeFormat })
}

func WithCSVDecoderOptionSetFieldValueFunc(f func(refrectType reflect.StructField, refrectValue reflect.Value, value string) (ok bool)) CSVDecoderOption {
	return csvDecoderOptionFunc(func(d *CSVDecoder) { d.setFieldValueFunc = f })
}
