package main

import (
	"io"
	"net/http"

	httpheaderschain "github.com/warrenhodg/opentracing-demo/tracing/chain/httpheaders"
	httpquerychain "github.com/warrenhodg/opentracing-demo/tracing/chain/httpquery"
	httpheadersmiddleware "github.com/warrenhodg/opentracing-demo/tracing/middleware/httpheaders"
	httpquerymiddleware "github.com/warrenhodg/opentracing-demo/tracing/middleware/httpquery"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// handleFoo passes the request onto another service
// to get the actual value. It passes the span via
// http headers to the other service
func handleFoo(w http.ResponseWriter, req *http.Request) {
	url := "http://localhost:8080/foo"
	req, _ := http.NewRequest("GET", url, nil)

	err := httpheaderschain.InjectSpan(req.Context(), req)
	if err != nil {
		// Log the error, but allow the request to proceed without
		// tracing having been setup
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.Header(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

// handleBar passes the request onto another service
// to get the actual value. It passes the span via
// http query to the other service
func handleFoo(w http.ResponseWriter, req *http.Request) {
	url := "http://localhost:8080/bar"
	req, _ := http.NewRequest("GET", url, nil)

	err := httpquerychain.InjectSpan(req.Context(), req, "telemetryspan")
	if err != nil {
		// Log the error, but allow the request to proceed without
		// tracing having been setup
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.Header(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func main() {
	cfg := jaegercfg.Configuration{
		ServiceName: "client2",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)

	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	http.DefaultServeMux.HandleFunc("/foo", httpheadersmiddleware.New(tracer, "foo", handleFoo).HandlerFunc)
	http.DefaultServeMux.HandleFunc("/bar", httpquerymiddleware.New(tracer, "bar", handleBar).HandlerFunc)

	http.ListenAndServe(":8081", nil)
}
