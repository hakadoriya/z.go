package tracez

import (
	"context"
	"path"
	"reflect"
	"testing"

	"github.com/hakadoriya/z.go/otelz/internal/consts"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestStart(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), (*trace.Tracer)(nil), noop.NewTracerProvider().Tracer(consts.TracerName))
		ctx, span := Start(ctx)
		defer span.End()

		if reflect.ValueOf(ctx).IsNil() {
			t.Errorf("❌: ctx: expected non-nil, actual nil")
		}

	})

	t.Run("error,", func(t *testing.T) {
		ctx := context.Background()
		ctx, span := Start(ctx)
		defer span.End()

		if reflect.ValueOf(ctx).IsNil() {
			t.Errorf("❌: ctx: expected non-nil, actual nil")
		}
	})

}

func Test_fullFuncName(t *testing.T) {
	t.Parallel()

	if expected, actual := path.Join(consts.ModuleName, "tracez.Test_fullFuncName"), wrapFullFuncName(); actual != expected {
		t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
	}
}

func wrapFullFuncName() string {
	return fullFuncName(1)
}
