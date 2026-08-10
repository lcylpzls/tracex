package tracex

import (
	"context"

	"github.com/lcylpzls/logx"
	"github.com/lcylpzls/tracex/internal/core"
	"go.opentelemetry.io/otel/attribute"
)

type (
	ExporterKind   = core.ExporterKind
	Config         = core.Config
	Manager        = core.Manager
	SpanSnapshot   = core.SpanSnapshot
	SpanEvent      = core.SpanEvent
	MemoryExporter = core.MemoryExporter
	TraceHook      = core.TraceHook
	TraceAttr      = core.TraceAttr
	LogHook        = core.LogHook
)

const (
	ExporterMemory   = core.ExporterMemory
	ExporterStdout   = core.ExporterStdout
	ExporterOTLPHTTP = core.ExporterOTLPHTTP
)

const (
	CodeInvalidConfig  = core.CodeInvalidConfig
	CodeExporterFailed = core.CodeExporterFailed
	CodeShutdownFailed = core.CodeShutdownFailed
)

func New(cfg Config) (*Manager, error) { return core.New(cfg) }
func NewHook(m *Manager) TraceHook     { return core.NewHook(m) }
func NewLogHook() *LogHook             { return core.NewLogHook() }
func WithBaggage(ctx context.Context, key, value string) context.Context {
	return core.WithBaggage(ctx, key, value)
}
func BaggageValue(ctx context.Context, key string) string { return core.BaggageValue(ctx, key) }
func NewMemoryExporter() *MemoryExporter                  { return core.NewMemoryExporter() }
func LogFields(ctx context.Context) logx.FieldGroup       { return core.LogFields(ctx) }
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	core.AddSpanEvent(ctx, name, attrs...)
}
func RecordError(ctx context.Context, err error) { core.RecordError(ctx, err) }
