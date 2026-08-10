package core

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// NewHook 返回标准的 TraceHook 实现：基于 Manager 创建 OpenTelemetry span。
// 可供任意家族库或第三方库通过各自的 WithTraceHook 选项接入，无需任何适配子包。
func NewHook(m *Manager) TraceHook {
	return &otelHook{m: m}
}

// otelHook 是 TraceHook 的 OpenTelemetry 标准实现。
type otelHook struct {
	m *Manager
}

// Start 在操作开始前创建 span，结束回调记录错误并结束 span。
func (h *otelHook) Start(ctx context.Context, name string, attrs ...TraceAttr) (context.Context, func(error)) {
	opts := make([]trace.SpanStartOption, 0, len(attrs))
	for _, a := range attrs {
		opts = append(opts, trace.WithAttributes(attribute.String(a.Key, a.Value)))
	}
	ctx, span := h.m.Start(ctx, name, opts...)
	return ctx, func(err error) {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}
}
