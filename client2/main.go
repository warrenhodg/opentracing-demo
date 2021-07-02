package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("health")
	defer func() {
		span.Finish()
	}()

	span.LogFields(
		log.String("uri", "health"),
		log.Int32("answertolifetheuniverseandeverything", 42))

	span.SetTag("relatedto", "health")

	url := "http://localhost:8080/health"
	req, _ := http.NewRequest("GET", url, nil)

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")

	tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func handleHealth2(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("health2")
	defer func() {
		span.Finish()
	}()

	span.SetTag("abc", "ABC")

	b := &strings.Builder{}
	tracer.Inject(span.Context(), opentracing.Binary, b)

	url := fmt.Sprintf("http://localhost:8080/health2?span=%s", url.QueryEscape(b.String()))
	req, _ := http.NewRequest("GET", url, nil)
	fmt.Printf("%v\n", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
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
			Param: 0.1,
			/*
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			*/
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

	http.ListenAndServe(":8081", nil)
}
