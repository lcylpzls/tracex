# tracex API 定版

> 版本：v0.1.0

## 1. Config

```go
type Config struct {
	ServiceName  string            // 服务名（必填）
	Version      string            // 服务版本
	Environment  string            // 部署环境
	SampleRatio  float64           // 采样率 0~1，0=默认 1
	BatchTimeout time.Duration     // 批量导出间隔，0=默认 5s
	Exporter     ExporterKind      // memory/stdout/otlp-http，空=stdout
	Writer       io.Writer         // stdout 输出目标
	OTLPEndpoint string            // OTLP/HTTP 端点
	OTLPInsecure bool              // 明文 HTTP
	OTLPHeaders  map[string]string // 附加请求头
	OTLPTimeout  time.Duration     // OTLP 请求超时（0=10s）
	SlowThreshold time.Duration    // 慢请求阈值（0=关闭）
	RouteNamer   func(*http.Request) string // 路由模板提取（可选）
	Sampler      sdktrace.Sampler  // 采样器（nil=采样率）
	SetGlobal    bool              // 注册为 OTel 全局组件
}
```

## 2. Manager

```go
func New(cfg Config) (*Manager, error)

func (m *Manager) Tracer() trace.Tracer
func (m *Manager) Propagator() propagation.TextMapPropagator
func (m *Manager) Config() Config
func (m *Manager) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
func (m *Manager) Inject(ctx context.Context, carrier propagation.TextMapCarrier)
func (m *Manager) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context
func (m *Manager) Middleware(next http.Handler) http.Handler
func (m *Manager) RoundTripper(next http.RoundTripper) http.RoundTripper
func (m *Manager) Spans() []SpanSnapshot   // 仅内存导出器
func (m *Manager) Shutdown(ctx context.Context) error
```

## 3. 日志联动

```go
func LogFields(ctx context.Context) logx.FieldGroup
```

## 4. Baggage 与事件

```go
func WithBaggage(ctx context.Context, key, value string) context.Context
func BaggageValue(ctx context.Context, key string) string
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue)
func RecordError(ctx context.Context, err error)

type LogHook struct{}
func NewLogHook() *LogHook // 注册到 logx.HookedLogger.AddHook
```

## 5. 内存导出器

```go
type SpanSnapshot struct {
	Name          string
	TraceID       string
	SpanID        string
	ParentSpanID  string
	Attributes    map[string]string
	StatusCode    string
	StatusMessage string
	Events        []SpanEvent
}

func NewMemoryExporter() *MemoryExporter
func (e *MemoryExporter) Spans() []SpanSnapshot
func (e *MemoryExporter) Reset()
```

## 6. 错误码

`tracex_invalid_config` / `tracex_exporter_failed` /
`tracex_shutdown_failed`（已注册 errx 分类，可用 `errx.Is` 匹配）。

## 7. webx 适配（adapters 子模块）

```go
import wx "github.com/lcylpzls/webx"
import txwebx "github.com/lcylpzls/tracex/adapters/webx"

s.UseGlobalMiddleware(txwebx.Middleware(m)) // m 为 *tracex.Manager
```

行为与标准中间件一致：链路提取/创建、路由级命名、状态码、
5xx 错误标记与慢请求事件。
