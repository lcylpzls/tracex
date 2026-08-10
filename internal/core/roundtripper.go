package core

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// RoundTripper 返回出站 HTTP 链路注入器：自动携带 traceparent、
// 记录客户端 span，并按状态码/错误标记结果。
func (m *Manager) RoundTripper(next http.RoundTripper) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &roundTripper{m: m, next: next}
}

// roundTripper 实现出站注入。
type roundTripper struct {
	m    *Manager
	next http.RoundTripper
}

// RoundTrip 注入链路并执行请求。
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := rt.m.Extract(req.Context(), propagation.HeaderCarrier(req.Header))
	spanName := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	ctx, span := rt.m.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("http.request.method", req.Method),
			attribute.String("url.full", req.URL.String()),
			attribute.String("server.address", req.URL.Host),
		),
	)
	defer span.End()
	rt.m.Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := rt.next.RoundTrip(req.WithContext(ctx))
	switch {
	case err != nil:
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	case resp != nil && resp.StatusCode >= http.StatusInternalServerError:
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", resp.StatusCode))
	default:
		span.SetStatus(codes.Ok, "")
	}
	if resp != nil {
		span.SetAttributes(attribute.Int("http.response.status_code", resp.StatusCode))
	}
	return resp, err
}
