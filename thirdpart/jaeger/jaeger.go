package jaeger

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

// NewJaegerTracer get tracer
func NewJaegerTracer(service string, endpoint string) opentracing.Tracer {
	sender := transport.NewHTTPTransport(endpoint)
	tracer, _ := jaeger.NewTracer(service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger)),
	)
	return tracer
}

// SpanNoParent
func SpanNoParent(tracer opentracing.Tracer, req *http.Request, name string, tag string, value string) {
	// 创建Span。
	span := tracer.StartSpan(name)
	// 设置Tag。
	span.SetTag(tag, value)
	// 将spanContext 透传到http 请求中
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	defer span.Finish()
}

// SpanParent
func SpanParent(parentSpan opentracing.Span, tracer opentracing.Tracer, req *http.Request, url string, tag string, value string) {
	// 创建Span。
	span := tracer.StartSpan(url, opentracing.ChildOf(parentSpan.Context()))
	// 设置Tag。
	span.SetTag(tag, value)
	// 将spanContext 透传到http 请求中
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	defer span.Finish()
}
