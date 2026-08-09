package tracex

import (
	"context"
	"sync"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SpanSnapshot Span 快照（内存导出器返回的只读视图）。
type SpanSnapshot struct {
	// Name Span 名称。
	Name string
	// TraceID 十六进制链路 ID。
	TraceID string
	// SpanID 十六进制 Span ID。
	SpanID string
	// ParentSpanID 父 Span ID（根 Span 为空）。
	ParentSpanID string
	// Attributes 属性快照。
	Attributes map[string]string
	// StatusCode 状态码（Unset/Ok/Error）。
	StatusCode string
	// StatusMessage 状态描述。
	StatusMessage string
}

// MemoryExporter 内存导出器：供测试与调试读取 Span。
type MemoryExporter struct {
	mu       sync.Mutex
	spans    []SpanSnapshot
	shutdown bool
}

// NewMemoryExporter 构造内存导出器。
func NewMemoryExporter() *MemoryExporter {
	return &MemoryExporter{}
}

// ExportSpans 接收并快照 Span。
func (e *MemoryExporter) ExportSpans(_ context.Context, spans []sdktrace.ReadOnlySpan) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.shutdown {
		return nil
	}
	for _, s := range spans {
		sc := s.SpanContext()
		snapshot := SpanSnapshot{
			Name:          s.Name(),
			TraceID:       sc.TraceID().String(),
			SpanID:        sc.SpanID().String(),
			ParentSpanID:  s.Parent().SpanID().String(),
			Attributes:    make(map[string]string, len(s.Attributes())),
			StatusCode:    s.Status().Code.String(),
			StatusMessage: s.Status().Description,
		}
		for _, kv := range s.Attributes() {
			snapshot.Attributes[string(kv.Key)] = kv.Value.String()
		}
		e.spans = append(e.spans, snapshot)
	}
	return nil
}

// Shutdown 关闭导出器（后续导出被忽略）。
func (e *MemoryExporter) Shutdown(context.Context) error {
	e.mu.Lock()
	e.shutdown = true
	e.mu.Unlock()
	return nil
}

// Spans 返回 Span 快照副本。
func (e *MemoryExporter) Spans() []SpanSnapshot {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]SpanSnapshot, len(e.spans))
	copy(out, e.spans)
	return out
}

// Reset 清空已收集的 Span。
func (e *MemoryExporter) Reset() {
	e.mu.Lock()
	e.spans = nil
	e.mu.Unlock()
}
