// Package jobx 提供 tracex 与 jobx 的链路追踪适配。
package jobx

import (
	"context"

	"github.com/lcylpzls/jobx"
	"github.com/lcylpzls/tracex"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Hook 实现 jobx.TraceHook。
type Hook struct {
	m *tracex.Manager
}

// NewHook 构造 jobx 追踪钩子。
func NewHook(m *tracex.Manager) *Hook {
	return &Hook{m: m}
}

// Start 在任务执行前启动 span，结束回调记录结果。
func (h *Hook) Start(ctx context.Context, name string, attrs ...jobx.TraceAttr) (context.Context, func(error)) {
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
