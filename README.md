# tracex

OpenTelemetry 链路追踪基座：统一 TracerProvider 管理、HTTP 中间件、
logx 日志字段与内存导出器，与 errx / logx 家族生态打通。

> 当前状态：**v1.2.0**。

## 特性

- TracerProvider 管理：stdout / OTLP/HTTP / 内存导出器；
- 标准 net/http 链路追踪中间件（W3C traceparent 透传）；
- 出站 HTTP 自动注入（RoundTripper 包装）；
- Baggage 便捷读写；
- 路由级 span 命名与慢请求标记；
- 日志钩子写入 span 事件、错误记录；
- 标准 `NewHook(m *Manager) TraceHook`：dbx/jobx/cachex/resiliencex/
  updatex/authx/winsvcx/filex 通过各自的 `WithTraceHook` 或
  `Config.TraceHook` 直接插拔接入，无需适配子包；
- webx 全局中间件直接使用标准中间件形态
  `s.UseGlobalMiddleware(m.Middleware)`；
- webx 全局中间件全请求覆盖（含 404/405），无追踪盲区。
- `LogFields(ctx)` 一键输出 logx trace_id/span_id 字段；
- 内存导出器（测试/调试友好，零外部依赖）；
- 统一 errx 错误码；简体中文注释与日志。

## 快速上手

```go
import (
	"context"
	"net/http"

	"github.com/lcylpzls/tracex"
)

func main() {
	m, err := tracex.New(tracex.Config{
		ServiceName: "order-service",
		Environment: "prod",
		Exporter:    tracex.ExporterOTLPHTTP,
		OTLPEndpoint: "collector:4318",
		OTLPInsecure: true,
	})
	if err != nil {
		panic(err)
	}
	defer m.Shutdown(context.Background())

	http.ListenAndServe(":8080", m.Middleware(http.DefaultServeMux))
}
```

日志联动：

```go
logger.Info("处理订单", tracex.LogFields(ctx)) // 自动带 trace_id/span_id
```

## 文档

- [架构详解](docs/architecture.md)
- [API 定版](docs/api.md)
- [集成示例](examples/README.md)
- [OTLP 采集部署](docs/otlp.md)
- [家族可观测性规范](docs/observability.md)

## License

MIT © [lcylpzls](https://github.com/lcylpzls)
