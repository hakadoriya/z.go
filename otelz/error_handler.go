package otelz

import "go.opentelemetry.io/otel"

var _ otel.ErrorHandler = (*ErrorHandleFunc)(nil)

type ErrorHandleFunc func(err error)

func (f ErrorHandleFunc) Handle(err error) {
	f(err)
}
