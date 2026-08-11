# tracex API 定版

> 当前版本：v1.3.0。公开 API 以 `go doc` 与 README 为准。

## 1. Config

```go
type Config struct {
	ServiceName   string            // 服务名（必填）
	Version       string            // 服务版本
	Environment   string            // 部署环境
	SampleRatio   float64           // 采样率 0~1，0=默认 1
	BatchTimeout  time.Duration     // 批量导出间隔，0=默认 5s
	Exporter      ExporterKind      // memory/stdout/otlp-http，空=stdout
	Writer        io.Writer         // stdout 输出目标
	OTLPEndpoint  string            // OTLP/HTTP 端点
	OTLPInsecure  bool              // 明文 HTTP
	OTLPHeaders   map[string]string // 附加请求头
	OTLPTimeout   time.Duration     // OTLP 请求超时（0=10s）
	SlowThreshold time.Duration     // 慢请求阈值（0=关闭）
	RouteNamer    func(*http.Request) string // 路由模板提取（可选）
	Sampler       sdktrace.Sampler  // 采样器（nil=按采样率）
	SetGlobal     bool              // 注册为 OTel 全局组件
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

## 3. 标准中间件（webx / net/http）

`Manager.Middleware` 返回标准 `func(http.Handler) http.Handler`，
可直接用于 webx 全局中间件或任意 `net/http` 服务：

```go
// webx：全局中间件覆盖 404/405 兜底请求，追踪无盲区
s.UseGlobalMiddleware(m.Middleware)

// 原生 net/http
http.ListenAndServe(":8080", m.Middleware(mux))
```

行为：链路提取/创建、URL 路由命名（可配 `RouteNamer`）、状态码记录、
5xx 错误标记与慢请求事件。

## 4. 家族插拔（标准 TraceHook）

`tracex.NewHook(m)` 返回家族统一的 `TraceHook`，各基座通过自身的
钩子选项直接接入，不需要任何适配子包：

```go
hook := tracex.NewHook(m)

dbx.Open(ctx, "mysql", dsn, dbx.WithTraceHook(hook))
jobx.NewDispatcher(jobx.WithTraceHook(hook))
cachex.New(cachex.WithTraceHook(hook))
resiliencex.NewCircuitBreaker(resiliencex.WithTraceHook(hook))
updatex.New(updatex.Config{TraceHook: hook, ...})
token.IssueRefreshToken(ctx, store, ttl, token.WithTraceHook(hook))
service.RunWithHook(name, hook) // winsvcx（Windows 平台）
filex.New(filex.Config{TraceHook: hook, ...})
```

httpx 出站通过传输层包装：

```go
client, _ := httpx.New(httpx.WithRoundTripperWrapper(m.RoundTripper))
```

## 5. 日志联动

```go
func LogFields(ctx context.Context) logx.FieldGroup
```

## 6. Baggage 与事件

```go
func WithBaggage(ctx context.Context, key, value string) context.Context
func BaggageValue(ctx context.Context, key string) string
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue)
func RecordError(ctx context.Context, err error)

type LogHook struct{}
func NewLogHook() *LogHook // 注册到 logx.HookedLogger.AddHook
```

## 7. 内存导出器

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

## 8. 错误码

`TRACEX_INVALID_CONFIG` / `TRACEX_EXPORTER_FAILED` /
`TRACEX_SHUTDOWN_FAILED`（已注册 errx 分类，可用 `errx.Is` 匹配）。
