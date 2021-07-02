package main

import (
	"fmt"
	"net/http"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := tracer.StartSpan("health", ext.RPCServerOption(spanCtx))
	span.LogFields(
		log.String("uri", "health"),
	)
	span.SetTag("relatedto", "health")
	span.Finish()
}

func handleHealth2(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	fmt.Printf("Query is %v\n", r.URL.String())
	s := r.URL.Query().Get("span")
	fmt.Printf("Span is %v\n", s)
	rdr := strings.NewReader(s)
	spanCtx, err := tracer.Extract(opentracing.Binary, rdr)
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	span := tracer.StartSpan("health2", ext.RPCServerOption(spanCtx))
	span.SetTag("def", "DEF")
	span.Finish()
}

func main() {
	cfg := jaegercfg.Configuration{
		ServiceName: "client1",
		Sampler: &jaegercfg.SamplerConfig{
			/*
				Type:  jaeger.SamplerTypeProbabilistic,
				Param: 0.1,
			*/
			Type:  jaeger.SamplerTypeConst,
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

	http.DefaultServeMux.HandleFunc("/health", handleHealth)
	http.DefaultServeMux.HandleFunc("/health2", handleHealth2)

	http.ListenAndServe(":8080", nil)
}
