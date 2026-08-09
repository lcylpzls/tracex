# OTLP 采集部署

## 1. Collector 配置（collector.yaml）

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
exporters:
  debug:
    verbosity: basic
service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [debug]
```

## 2. Docker 启动

```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    command: ["--config=/etc/otelcol/config.yaml"]
    volumes:
      - ./collector.yaml:/etc/otelcol/config.yaml
    ports:
      - "4318:4318"
```

## 3. tracex 接入

```go
m, err := tracex.New(tracex.Config{
	ServiceName:   "order-service",
	Version:       "1.0.0",
	Environment:   "prod",
	Exporter:      tracex.ExporterOTLPHTTP,
	OTLPEndpoint:  "collector:4318",
	OTLPInsecure:  true,
	OTLPHeaders:   map[string]string{"Authorization": "Bearer <采集令牌>"},
	OTLPTimeout:   10 * time.Second,
	SlowThreshold: 500 * time.Millisecond,
	RouteNamer: func(r *http.Request) string {
		return r.URL.Path // 框架适配时返回路由模板
	},
})
```

## 4. 生产建议

- Collector 使用 HTTPS 与令牌认证，`OTLPInsecure` 仅限内网；
- 采样率按流量设置 `SampleRatio`（如 0.1）；
- 日志通过 `LogFields(ctx)` / `LogHook` 关联 trace_id/span_id；
- 服务下线前 `defer m.Shutdown(ctx)` 确保 span 刷新导出。
