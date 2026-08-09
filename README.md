# tracex

OpenTelemetry 链路追踪基座：统一 TracerProvider 管理、HTTP 中间件、
logx 日志字段与内存导出器，与 errx / logx 家族生态打通。

> 当前状态：**v0.6.0（自主打磨）**。

## 特性

- TracerProvider 管理：stdout / OTLP/HTTP / 内存导出器；
- 标准 net/http 链路追踪中间件（W3C traceparent 透传）；
- 出站 HTTP 自动注入（RoundTripper 包装）；
- Baggage 便捷读写；
- 路由级 span 命名与慢请求标记；
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

- [设计定版](docs/design.md)
- [架构详解](docs/architecture.md)
- [API 定版](docs/api.md)
- [版本路线](docs/roadmap.md)

## License

MIT © [lcylpzls](https://github.com/lcylpzls)
