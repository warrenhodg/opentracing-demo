package httpheaders

import (
	"context"
	"net/http"

	"github.com/warrenhodg/opentracing-demo/tracing"

	opentracing "github.com/opentracing/opentracing-go"
)

// InjectSpan injects a span extracted from the context into
// the url of the given request object
func InjectSpan(ctx context.Context, req *http.Request) error {
	tracer, span, err := tracing.GetTracerAndSpan(ctx)
	if err != nil {
		return err
	}

	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	return nil
}
