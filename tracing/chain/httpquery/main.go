package httpquery

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/warrenhodg/opentracing-demo/tracing"

	opentracing "github.com/opentracing/opentracing-go"
)

// InjectSpan injects a span extracted from the context into
// the url of the given request object
func InjectSpan(ctx context.Context, req *http.Request, paramName string) error {
	tracer, span, err := tracing.GetTracerAndSpan(ctx)
	if err != nil {
		return err
	}

	b := &strings.Builder{}
	tracer.Inject(span.Context(), opentracing.Binary, b)

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return err
	}
	query.Add(paramName, b.String())

	req.URL.RawQuery = query.Encode()

	return nil
}
