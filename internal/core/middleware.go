package core

import (
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Middleware 返回标准 net/http 链路追踪中间件：
// 提取入站链路上下文、记录请求属性与状态码，5xx 标记错误。
func (m *Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := m.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := m.Start(ctx, r.URL.Path,
			trace.WithAttributes(
				attribute.String("http.request.method", r.Method),
				attribute.String("url.path", r.URL.Path),
				attribute.String("url.scheme", r.URL.Scheme),
				attribute.String("server.address", r.Host),
				attribute.String("user_agent.original", r.UserAgent()),
			),
		)
		defer span.End()
		start := time.Now()
		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(ctx))
		span.SetAttributes(attribute.Int("http.response.status_code", rw.status))
		// 路由级命名：框架适配通过 RouteNamer 提供路由模板。
		if m.cfg.RouteNamer != nil {
			if route := m.cfg.RouteNamer(r); route != "" {
				span.SetName(r.Method + " " + route)
			}
		}
		elapsed := time.Since(start)
		if m.cfg.SlowThreshold > 0 && elapsed > m.cfg.SlowThreshold {
			span.SetAttributes(attribute.Int64("request.duration_ms", elapsed.Milliseconds()))
			span.AddEvent("slow", trace.WithAttributes(
				attribute.Int64("elapsed_ms", elapsed.Milliseconds()),
			))
		}
		if rw.status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", rw.status))
		}
	})
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
