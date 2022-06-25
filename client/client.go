package main

import (
	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"net/http"
)

func main() {
	tracer := opentracing.GlobalTracer()
	cfg := &config.Configuration{
		ServiceName: "client",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	//1
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	defer closer.Close()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	//2
	clientSpan := tracer.StartSpan("clientspan")
	defer clientSpan.Finish()

	url := "http://localhost:8080/publish"
	req, _ := http.NewRequest("GET", url, nil)

	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, url)
	ext.HTTPMethod.Set(clientSpan, "GET")

	tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println(resp)
}
