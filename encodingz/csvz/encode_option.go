package csvz

import "encoding/csv"

type (
	CSVEncoderOption     interface{ apply(d *CSVEncoder) }
	csvEncoderOptionFunc func(d *CSVEncoder)
)

func (f csvEncoderOptionFunc) apply(d *CSVEncoder) { f(d) }

func WithCSVEncoderOptionCSVWriter(w *csv.Writer) CSVEncoderOption {
	return csvEncoderOptionFunc(func(d *CSVEncoder) { d.w = w })
}

func WithCSVEncoderOptionTagName(tagName string) CSVEncoderOption {
	return csvEncoderOptionFunc(func(d *CSVEncoder) { d.tagName = tagName })
}

func WithCSVEncoderOptionTimeFormat(timeFormat string) CSVEncoderOption {
	return csvEncoderOptionFunc(func(d *CSVEncoder) { d.timeFormat = timeFormat })
}
