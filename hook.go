package tracex

import (
	"github.com/lcylpzls/logx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// LogHook logx 日志钩子：把日志写入当前 span 的 log 事件，
// 实现日志与链路事件关联。
type LogHook struct{}

// NewLogHook 构造日志钩子。
func NewLogHook() *LogHook {
	return &LogHook{}
}

// OnLog 在日志写入后调用（logx 异步触发）。
func (h *LogHook) OnLog(e *logx.Entry) {
	span := trace.SpanFromContext(e.Context())
	if !span.IsRecording() {
		return
	}
	span.AddEvent("log", trace.WithAttributes(
		attribute.String("log.level", e.Level.String()),
		attribute.String("log.message", e.Message),
	))
}
