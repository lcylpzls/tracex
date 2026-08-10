// Package tracex 提供 OpenTelemetry 链路追踪基座：
// TracerProvider 管理、HTTP 中间件、logx 日志字段与内存导出器。
package core

import (
	"context"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/lcylpzls/errx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// 可替换的系统构造（测试注入用）。
var (
	buildStdout = func(w io.Writer) (sdktrace.SpanExporter, error) {
		return stdouttrace.New(stdouttrace.WithWriter(w))
	}
	buildOTLP = func(ctx context.Context, opts ...otlptracehttp.Option) (sdktrace.SpanExporter, error) {
		return otlptracehttp.New(ctx, opts...)
	}
	shutdownProvider = func(ctx context.Context, p *sdktrace.TracerProvider) error {
		return p.Shutdown(ctx)
	}
)

// ExporterKind 导出器类型。
type ExporterKind string

const (
	// ExporterMemory 内存导出器（测试/调试）。
	ExporterMemory ExporterKind = "memory"
	// ExporterStdout 标准输出导出器（默认）。
	ExporterStdout ExporterKind = "stdout"
	// ExporterOTLPHTTP OTLP/HTTP 导出器（生产采集）。
	ExporterOTLPHTTP ExporterKind = "otlp-http"
)

// Config 追踪器配置。
type Config struct {
	// ServiceName 服务名（必填）。
	ServiceName string
	// Version 服务版本（可选）。
	Version string
	// Environment 部署环境（可选）。
	Environment string
	// SampleRatio 采样率 0~1，0 使用默认 1（全采样）。
	SampleRatio float64
	// BatchTimeout 批量导出间隔，0 使用默认 5s。
	BatchTimeout time.Duration
	// Exporter 导出器类型，空使用 stdout。
	Exporter ExporterKind
	// Writer stdout 导出器输出目标，nil 使用 os.Stdout。
	Writer io.Writer
	// OTLPEndpoint OTLP/HTTP 端点（host:port）。
	OTLPEndpoint string
	// OTLPInsecure 是否使用明文 HTTP。
	OTLPInsecure bool
	// OTLPHeaders OTLP/HTTP 附加请求头。
	OTLPHeaders map[string]string
	// OTLPTimeout OTLP/HTTP 请求超时，0 使用默认 10s。
	OTLPTimeout time.Duration
	// SlowThreshold 慢请求阈值；超过时记录 slow 事件，0 关闭。
	SlowThreshold time.Duration
	// RouteNamer 可选：从请求中提取路由模板（如 /users/{id}）用于
	// span 命名；返回空串时保持默认命名。
	RouteNamer func(r *http.Request) string
	// Sampler 可选采样器；nil 使用 TraceIDRatioBased(SampleRatio)。
	Sampler sdktrace.Sampler
	// SetGlobal 是否注册为 OTel 全局 TracerProvider 与传播器。
	SetGlobal bool
}

// Manager 追踪管理器：持有 TracerProvider、导出器与传播器。
type Manager struct {
	cfg        Config
	provider   *sdktrace.TracerProvider
	exporter   sdktrace.SpanExporter
	mem        *MemoryExporter
	propagator propagation.TextMapPropagator
	tracer     trace.Tracer
	mu         sync.Mutex
	closed     bool
}

// New 构造追踪管理器。
func New(cfg Config) (*Manager, error) {
	if cfg.ServiceName == "" {
		return nil, errx.NewCode(CodeInvalidConfig, "服务名不能为空")
	}
	if cfg.SampleRatio < 0 || cfg.SampleRatio > 1 {
		return nil, errx.NewCode(CodeInvalidConfig, "采样率必须在 0~1 之间")
	}
	if cfg.SampleRatio == 0 {
		cfg.SampleRatio = 1
	}
	if cfg.BatchTimeout <= 0 {
		cfg.BatchTimeout = 5 * time.Second
	}
	exp, mem, err := buildExporter(cfg)
	if err != nil {
		return nil, err
	}
	attrs := []attribute.KeyValue{
		attribute.String("service.name", cfg.ServiceName),
	}
	if cfg.Version != "" {
		attrs = append(attrs, attribute.String("service.version", cfg.Version))
	}
	if cfg.Environment != "" {
		attrs = append(attrs, attribute.String("deployment.environment", cfg.Environment))
	}
	sampler := cfg.Sampler
	if sampler == nil {
		sampler = sdktrace.TraceIDRatioBased(cfg.SampleRatio)
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(cfg.BatchTimeout)),
		sdktrace.WithResource(sdkresource.NewSchemaless(attrs...)),
		sdktrace.WithSampler(sampler),
	)
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	if cfg.SetGlobal {
		otel.SetTracerProvider(provider)
		otel.SetTextMapPropagator(propagator)
	}
	m := &Manager{
		cfg:        cfg,
		provider:   provider,
		exporter:   exp,
		mem:        mem,
		propagator: propagator,
		tracer:     provider.Tracer("github.com/lcylpzls/tracex"),
	}
	return m, nil
}

// buildExporter 按配置构造导出器。
func buildExporter(cfg Config) (sdktrace.SpanExporter, *MemoryExporter, error) {
	switch cfg.Exporter {
	case "", ExporterStdout:
		w := cfg.Writer
		if w == nil {
			w = os.Stdout
		}
		exp, err := buildStdout(w)
		if err != nil {
			return nil, nil, errx.WrapCode(err, CodeExporterFailed, "创建 stdout 导出器失败")
		}
		return exp, nil, nil
	case ExporterMemory:
		mem := NewMemoryExporter()
		return mem, mem, nil
	case ExporterOTLPHTTP:
		if cfg.OTLPEndpoint == "" {
			return nil, nil, errx.NewCode(CodeInvalidConfig, "OTLP/HTTP 端点不能为空")
		}
		opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.OTLPEndpoint)}
		if cfg.OTLPInsecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		if len(cfg.OTLPHeaders) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(cfg.OTLPHeaders))
		}
		timeout := cfg.OTLPTimeout
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		opts = append(opts, otlptracehttp.WithTimeout(timeout))
		exp, err := buildOTLP(context.Background(), opts...)
		if err != nil {
			return nil, nil, errx.WrapCode(err, CodeExporterFailed, "创建 OTLP/HTTP 导出器失败")
		}
		return exp, nil, nil
	default:
		return nil, nil, errx.NewCode(CodeInvalidConfig, "不支持的导出器类型："+string(cfg.Exporter))
	}
}

// Tracer 返回管理器持有的 Tracer。
func (m *Manager) Tracer() trace.Tracer {
	return m.tracer
}

// Propagator 返回组合传播器（TraceContext + Baggage）。
func (m *Manager) Propagator() propagation.TextMapPropagator {
	return m.propagator
}

// Config 返回配置副本。
func (m *Manager) Config() Config {
	return m.cfg
}

// Start 启动一个 Span（便捷封装）。
func (m *Manager) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return m.tracer.Start(ctx, spanName, opts...)
}

// Inject 将链路上下文写入 carrier（出站请求）。
func (m *Manager) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	m.propagator.Inject(ctx, carrier)
}

// Extract 从 carrier 提取链路上下文（入站请求）。
func (m *Manager) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return m.propagator.Extract(ctx, carrier)
}

// Spans 返回内存导出器中的 Span 快照（非内存导出器返回 nil）。
func (m *Manager) Spans() []SpanSnapshot {
	if m.mem == nil {
		return nil
	}
	return m.mem.Spans()
}

// Shutdown 刷新并关闭追踪器（幂等）。
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	m.mu.Unlock()
	if err := shutdownProvider(ctx, m.provider); err != nil {
		return errx.WrapCode(err, CodeShutdownFailed, "关闭追踪器失败")
	}
	return nil
}
