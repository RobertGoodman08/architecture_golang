package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"

	log "architecture_go/pkg/type/logger"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		if span == nil {
			span = StartSpanWithHeader(&c.Request.Header, "rest-request-"+c.Request.Method, c.Request.Method, c.Request.URL.Path)
		}
		defer span.Finish()
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

		if traceID, ok := span.Context().(jaeger.SpanContext); ok {
			c.Header("uber-trace-id", traceID.TraceID().String())
		}

		c.Next()

		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))

		if len(c.Errors) == 0 {
			log.Info("", getContextFields(c)...)
		}
	}
}

func StartSpanWithHeader(header *http.Header, operationName, method, path string) opentracing.Span {
	var wireContext opentracing.SpanContext

	if header != nil {
		wireContext, _ = opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*header))
	}

	return StartSpanWithParent(wireContext, operationName, method, path)
}

// StartSpanWithParent will start a new span with a parent span.
// example:
//      span:= StartSpanWithParent(c.Get("tracing-context"),
func StartSpanWithParent(parent opentracing.SpanContext, operationName, method, path string) opentracing.Span {
	options := []opentracing.StartSpanOption{
		opentracing.Tag{Key: ext.SpanKindRPCServer.Key, Value: ext.SpanKindRPCServer.Value},
		opentracing.Tag{Key: string(ext.HTTPMethod), Value: method},
		opentracing.Tag{Key: string(ext.HTTPUrl), Value: path},
	}
	if parent != nil {
		options = append(options, opentracing.ChildOf(parent))
	}

	return opentracing.StartSpan(operationName, options...)
}
