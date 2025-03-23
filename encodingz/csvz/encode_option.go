package csvz

import "encoding/csv"

type (
	CSVEncoderOption     interface{ apply(d *CSVEncoder) }
	csvEncoderOptionFunc func(d *CSVEncoder)
)

func (f csvEncoderOptionFunc) apply(d *CSVEncoder) { f(d) }

func WithCSVEncoderOptionCSVWriterModifier(modifier func(w *csv.Writer) *csv.Writer) CSVEncoderOption {
	return csvEncoderOptionFunc(func(d *CSVEncoder) { d.csvWriter = modifier(d.csvWriter) })
}

func WithCSVEncoderOptionTagName(tagName string) CSVEncoderOption {
	return csvEncoderOptionFunc(func(d *CSVEncoder) { d.tagName = tagName })
}

func WithCSVEncoderOptionTimeFormat(timeFormat string) CSVEncoderOption {
	return csvEncoderOptionFunc(func(d *CSVEncoder) { d.timeFormat = timeFormat })
}
