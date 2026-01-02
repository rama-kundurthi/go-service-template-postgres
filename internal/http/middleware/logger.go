package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		requestID, _ := c.Get("request_id")
		span := trace.SpanFromContext(c.Request.Context())
		sc := span.SpanContext()

		attrs := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"request_id", requestID,
		}

		if sc.IsValid() {
			attrs = append(attrs,
				"trace_id", sc.TraceID().String(),
				"span_id", sc.SpanID().String(),
			)
		}

		// Gin collects errors in c.Errors
		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		// Use Info for all; you can bump to Warn/Error by status if desired.
		logger.Info("http_request", attrs...)
	}
}
