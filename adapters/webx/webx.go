// Package webx 提供 tracex 与 webx 框架的中间件适配。
package webx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lcylpzls/tracex"
	wx "github.com/lcylpzls/webx/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Middleware 返回 webx 全局中间件：提取/创建链路、记录请求属性、
// 路由级命名、状态码与慢请求标记。
func Middleware(m *tracex.Manager) wx.HandlerFunc {
	return func(c *wx.Context) {
		req := c.Request()
		ctx := m.Extract(req.Context(), propagation.HeaderCarrier(req.Header))
		ctx, span := m.Start(ctx, req.URL.Path,
			trace.WithAttributes(
				attribute.String("http.request.method", req.Method),
				attribute.String("url.path", req.URL.Path),
				attribute.String("url.scheme", req.URL.Scheme),
				attribute.String("server.address", req.Host),
				attribute.String("user_agent.original", req.UserAgent()),
			),
		)
		defer span.End()
		start := time.Now()
		rw := &statusWriter{ResponseWriter: c.Writer(), status: http.StatusOK}
		c.SetRequest(req.WithContext(ctx))
		c.SetWriter(rw)
		c.Next()
		span.SetAttributes(attribute.Int("http.response.status_code", rw.status))
		if route := c.Route(); route != "" {
			span.SetName(req.Method + " " + route)
		}
		if slow := m.Config().SlowThreshold; slow > 0 {
			if elapsed := time.Since(start); elapsed > slow {
				span.SetAttributes(attribute.Int64("request.duration_ms", elapsed.Milliseconds()))
				span.AddEvent("slow", trace.WithAttributes(
					attribute.Int64("elapsed_ms", elapsed.Milliseconds()),
				))
			}
		}
		if rw.status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", rw.status))
		}
	}
}

// statusWriter 捕获响应状态码。
type statusWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader 记录状态码并透传。
func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
