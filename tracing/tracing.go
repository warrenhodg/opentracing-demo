package tracing

import (
	"context"
	"errors"

	opentracing "github.com/opentracing/opentracing-go"
)

// ContextKey type is used to store Tracing information
// in context.Context instances.
type ContextKey = struct {
	value string
}

// ContextKeyTracer is the key used to store the
// actual opentelemetry tracer object in a
// context.Context instance
var contextKeyTracer = &ContextKey{"tracer"}

// ContextKeySpan is the key used to store the
// actual opentelemetry span object in a
// context.Context instance
var contextKeySpan = &ContextKey{"span"}

// DefaultTelemetryParam is the default name of parameter
// used to transmit telemetry span information.
// It is provided here as a means of standardization where
// applicable.
// For example, it could be a url query parameter.
// A different parameter name can be used, but this
// information must be communicated to teams developing
// related services.
const DefaultTelemetryParam = "telemetrycontext"

// SetTracerAndSpan sets the tracer and span values in the context, and returns a new context
func SetTracerAndSpan(ctx context.Context, tracer opentracing.Tracer, span opentracing.Span) context.Context {
	ctx = context.WithValue(ctx, contextKeyTracer, tracer)
	ctx = context.WithValue(ctx, contextKeySpan, span)
	return ctx
}

// GetTracerAndSpan gets the tracer and span values from the context.
// Return values are guaranteed to be non-nil interfaces or an error
func GetTracerAndSpan(ctx context.Context) (tracer opentracing.Tracer, span opentracing.Span, err error) {
	var ok bool

	tracerInterface := ctx.Value(contextKeyTracer)
	if tracerInterface == nil {
		return nil, nil, errors.New("tracer not set in context")
	}
	if tracer, ok = tracerInterface.(opentracing.Tracer); !ok {
		return nil, nil, errors.New("tracer not valid in context")
	}

	spanInterface := ctx.Value(contextKeySpan)
	if spanInterface == nil {
		return nil, nil, errors.New("span not set in context")
	}
	if span, ok = spanInterface.(opentracing.Span); !ok {
		return nil, nil, errors.New("span not valid in context")
	}

	return tracer, span, nil
}
