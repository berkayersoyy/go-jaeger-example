package main

import (
	"errors"
	"fmt"
	"github.com/berkayersoyy/go-jaeger-example/metric"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	metricsMiddleware := metric.NewMetricsMiddleware()

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/publish", publishHandler).Methods(http.MethodGet)

	r.Use(metricsMiddleware.Metrics)

	http.ListenAndServe(":8080", r)
}
func publishHandler(w http.ResponseWriter, r *http.Request) {
	cfg := jaegercfg.Configuration{
		ServiceName: "service_name",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	serverSpan := tracer.StartSpan("main.Publish", ext.RPCServerOption(spanCtx))
	err = errors.New("error from server side")
	serverSpan.LogKV("Error", err)
	w.WriteHeader(200)
	//w.Write([]byte("Publish"))
	defer serverSpan.Finish()
}
