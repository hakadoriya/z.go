package tracez

import (
	"context"
	"runtime"

	"github.com/hakadoriya/z.go/otelz/internal/consts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Start(ctx context.Context, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	//nolint:ireturn,spancheck
	return tracerFromContext(ctx).Start(ctx, fullFuncName(1), opts...)
}

func StartWithSpanName(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	//nolint:ireturn,spancheck
	return tracerFromContext(ctx).Start(ctx, spanName, opts...)
}

func StartWithSpanNameSuffix(ctx context.Context, spanNameSuffix string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	//nolint:ireturn,spancheck
	return tracerFromContext(ctx).Start(ctx, fullFuncName(1)+"."+spanNameSuffix, opts...)
}

func StartFuncWithSpanName(ctx context.Context, spanName string, f func(ctx context.Context) (err error), opts ...trace.SpanStartOption) error {
	ctx, span := tracerFromContext(ctx).Start(ctx, spanName, opts...)
	defer span.End()
	return f(ctx)
}

func StartFuncWithSpanNameSuffix(ctx context.Context, spanNameSuffix string, f func(ctx context.Context) (err error), opts ...trace.SpanStartOption) error {
	ctx, span := tracerFromContext(ctx).Start(ctx, fullFuncName(1)+"."+spanNameSuffix, opts...)
	defer span.End()
	return f(ctx)
}

func tracerFromContext(ctx context.Context) trace.Tracer {
	tracer, ok := ctx.Value((*trace.Tracer)(nil)).(trace.Tracer)
	if !ok {
		return otel.GetTracerProvider().Tracer(consts.TracerName)
	}

	return tracer
}

func fullFuncName(skip int) (funcName string) {
	pc, _, _, _ := runtime.Caller(skip + 1) //nolint:dogsled
	return runtime.FuncForPC(pc).Name()
}
