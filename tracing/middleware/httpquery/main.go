package httpquery

import (
	"net/http"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/warrenhodg/opentracing-demo/tracing"
)

// Middleware wrap's a request in an
// tracing span, potentially continuing
// the span passed in via http query from
// the requestor
type Middleware struct {
	tracer        opentracing.Tracer
	operationName string
	queryParam    string
	handler       func(http.ResponseWriter, *http.Request)
}

// New instantiates a new Middleware instance
func New(tracer opentracing.Tracer, operationName string, queryParam string, handler func(http.ResponseWriter, *http.Request)) *Middleware {
	return &Middleware{
		tracer:        tracer,
		operationName: operationName,
		queryParam:    queryParam,
		handler:       handler,
	}
}

// HandlerFunc wraps the call in a span, and continues then passes the request
// to another function handler
func (m *Middleware) HandlerFunc(w http.ResponseWriter, req *http.Request) {
	var (
		spanCtx opentracing.SpanContext
		err     error
	)

	if m.queryParam != "" {
		v := req.URL.Query().Get(m.queryParam)
		if v != "" {
			r := strings.NewReader(v)
			spanCtx, err = m.tracer.Extract(opentracing.Binary, r)
			if err != nil {
				// Log the error. If spanCtx is nil here, then we are
				// unable to continue an existing span contact,
				// and will simply create a new one
			}
		}
	}

	// if spanCtx is nil here, then we are unable to continue an existing span contact,
	// and will simply create a new one
	span := m.tracer.StartSpan(m.operationName, ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx := tracing.SetTracerAndSpan(req.Context(), m.tracer, span)

	req = req.WithContext(ctx)
	m.handler(w, req)
}
