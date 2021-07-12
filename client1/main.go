package main

import (
	"net/http"

	"github.com/warrenhodg/opentracing-demo/tracing"
	httpheadersmiddleware "github.com/warrenhodg/opentracing-demo/tracing/middleware/httpheaders"
	httpquerymiddleware "github.com/warrenhodg/opentracing-demo/tracing/middleware/httpquery"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func handleFoo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleBar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	cfg := jaegercfg.Configuration{
		ServiceName: "client1",
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
	http.DefaultServeMux.HandleFunc("/bar", httpquerymiddleware.New(tracer, "bar", tracing.DefaultTelemetryParam, handleBar).HandlerFunc)

	http.ListenAndServe(":8080", nil)
}
