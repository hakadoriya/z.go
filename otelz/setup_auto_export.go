package otelz

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func SetupAutoExport(ctx context.Context, opts ...AutoExportOption) (shutdown func(ctx context.Context) (err error), err error) {
	c := new(autoexportConfig)
	for _, o := range opts {
		o.apply(c)
	}

	var (
		shutdownFuncs []func(context.Context) error
	)
	shutdown = func(ctx context.Context) (err error) {
		for _, fn := range reverse(shutdownFuncs) {
			if e := fn(ctx); e != nil {
				err = errors.Join(err, e)
			}
		}
		shutdownFuncs = nil
		return err
	}

	// resource
	res, err := resource.New(ctx, deepDistinct(c.resourceOptions)...)
	if err != nil {
		return shutdown, fmt.Errorf("resource.New: %w", err)
	}

	// propagator
	propagator := autoprop.NewTextMapPropagator(deepDistinct(append([]propagation.TextMapPropagator{propagation.TraceContext{}, propagation.Baggage{}}, c.textMapPropagators...))...)
	otel.SetTextMapPropagator(propagator)

	// trace
	spanExporter := c.spanExporter
	if (spanExporter == nil) || reflect.ValueOf(spanExporter).IsNil() {
		spanExporter, err = autoexport.NewSpanExporter(ctx, deepDistinct(c.spanExporterOptions)...)
		if err != nil {
			return shutdown, fmt.Errorf("autoexport.NewSpanExporter: %w", err)
		}
	}
	shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
		if err := spanExporter.Shutdown(ctx); err != nil {
			return fmt.Errorf("spanExporter.Shutdown: %w", err)
		}
		return nil
	})
	tracerProvider := trace.NewTracerProvider(deepDistinct(append([]trace.TracerProviderOption{trace.WithBatcher(spanExporter), trace.WithResource(res)}, c.tracerProviderOptions...))...)
	shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
		if err := tracerProvider.ForceFlush(ctx); err != nil {
			return fmt.Errorf("tracerProvider.ForceFlush: %w", err)
		}
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("tracerProvider.Shutdown: %w", err)
		}
		return nil
	})
	otel.SetTracerProvider(tracerProvider)

	// metric
	metricReader := c.metricReader
	if (metricReader == nil) || reflect.ValueOf(metricReader).IsNil() {
		metricReader, err = autoexport.NewMetricReader(ctx, deepDistinct(c.metricReaderOptions)...)
		if err != nil {
			return shutdown, fmt.Errorf("autoexport.NewMetricReader: %w", err)
		}
	}
	shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
		if err := metricReader.Shutdown(ctx); err != nil {
			return fmt.Errorf("metricReader.Shutdown: %w", err)
		}
		return nil
	})
	meterProvider := metric.NewMeterProvider(deepDistinct(append([]metric.Option{metric.WithReader(metricReader), metric.WithResource(res)}, c.metricProviderOptions...))...)
	shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
		if err := meterProvider.ForceFlush(ctx); err != nil {
			return fmt.Errorf("meterProvider.ForceFlush: %w", err)
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("meterProvider.Shutdown: %w", err)
		}
		return nil
	})
	otel.SetMeterProvider(meterProvider)

	return shutdown, nil
}

type AutoExportOption interface {
	apply(*autoexportConfig)
}

type autoexportConfig struct {
	// resource
	resourceOptions []resource.Option
	// propagator
	textMapPropagators []propagation.TextMapPropagator
	// trace
	spanExporter          trace.SpanExporter
	spanExporterOptions   []autoexport.SpanOption
	tracerProviderOptions []trace.TracerProviderOption
	// metric
	metricReader          metric.Reader
	metricReaderOptions   []autoexport.MetricOption
	metricProviderOptions []metric.Option
}

// resource

func WithResourceOptions(opts ...resource.Option) AutoExportOption {
	return &withAutoExportResourceOptions{opts: opts}
}

type withAutoExportResourceOptions struct {
	opts []resource.Option
}

func (w *withAutoExportResourceOptions) apply(c *autoexportConfig) {
	c.resourceOptions = append(c.resourceOptions, w.opts...)
}

// propagator

func WithAutoExportTextMapPropagators(propagators ...propagation.TextMapPropagator) AutoExportOption {
	return &withAutoExportTextMapPropagators{propagators: propagators}
}

type withAutoExportTextMapPropagators struct {
	propagators []propagation.TextMapPropagator
}

func (w *withAutoExportTextMapPropagators) apply(c *autoexportConfig) {
	c.textMapPropagators = append(c.textMapPropagators, w.propagators...)
}

// trace

func WithAutoExportSpanExporter(exporter trace.SpanExporter) AutoExportOption {
	return &withAutoExportSpanExporter{exporter: exporter}
}

type withAutoExportSpanExporter struct {
	exporter trace.SpanExporter
}

func (w *withAutoExportSpanExporter) apply(c *autoexportConfig) {
	c.spanExporter = w.exporter
}

func WithAutoExportSpanExporterOptions(opts ...autoexport.SpanOption) AutoExportOption {
	return &withAutoExportSpanExporterOptions{opts: opts}
}

type withAutoExportSpanExporterOptions struct {
	opts []autoexport.SpanOption
}

func (w *withAutoExportSpanExporterOptions) apply(c *autoexportConfig) {
	c.spanExporterOptions = append(c.spanExporterOptions, w.opts...)
}

func WithAutoExportTracerProviderOptions(opts ...trace.TracerProviderOption) AutoExportOption {
	return &withAutoExportTracerProviderOptions{opts: opts}
}

type withAutoExportTracerProviderOptions struct {
	opts []trace.TracerProviderOption
}

func (w *withAutoExportTracerProviderOptions) apply(c *autoexportConfig) {
	c.tracerProviderOptions = append(c.tracerProviderOptions, w.opts...)
}

// metric

func WithAutoExportMetricReader(reader metric.Reader) AutoExportOption {
	return &withAutoExportMetricReader{reader: reader}
}

type withAutoExportMetricReader struct {
	reader metric.Reader
}

func (w *withAutoExportMetricReader) apply(c *autoexportConfig) {
	c.metricReader = w.reader
}

func WithAutoExportMetricReaderOptions(opts ...autoexport.MetricOption) AutoExportOption {
	return &withAutoExportMetricReaderOptions{opts: opts}
}

type withAutoExportMetricReaderOptions struct {
	opts []autoexport.MetricOption
}

func (w *withAutoExportMetricReaderOptions) apply(c *autoexportConfig) {
	c.metricReaderOptions = append(c.metricReaderOptions, w.opts...)
}

func WithAutoExportMetricProviderOptions(opts ...metric.Option) AutoExportOption {
	return &withAutoExportMetricProviderOptions{opts: opts}
}

type withAutoExportMetricProviderOptions struct {
	opts []metric.Option
}

func (w *withAutoExportMetricProviderOptions) apply(c *autoexportConfig) {
	c.metricProviderOptions = append(c.metricProviderOptions, w.opts...)
}
