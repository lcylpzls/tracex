package tracex

import (
	"context"

	"github.com/lcylpzls/logx"
	"go.opentelemetry.io/otel/trace"
)

// LogFields 从上下文提取链路信息生成 logx 字段；
// 无有效链路时返回空字段组。
func LogFields(ctx context.Context) logx.FieldGroup {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return logx.Fields()
	}
	return logx.Fields(
		logx.String("trace_id", sc.TraceID().String()),
		logx.String("span_id", sc.SpanID().String()),
	)
}
