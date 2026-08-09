# tracex 集成示例

## 服务端（入站）

```go
m, err := tracex.New(tracex.Config{
	ServiceName:  "order-service",
	Environment:  "prod",
	Exporter:     tracex.ExporterOTLPHTTP,
	OTLPEndpoint: "collector:4318",
	OTLPInsecure: true,
	RouteNamer: func(r *http.Request) string {
		// 框架适配：返回路由模板，如 webx/gin 路由名。
		return r.URL.Path
	},
})
if err != nil {
	panic(err)
}
defer m.Shutdown(context.Background())

http.ListenAndServe(":8080", m.Middleware(mux))
```

## 客户端（出站）

```go
client := &http.Client{Transport: m.RoundTripper(http.DefaultTransport)}
resp, err := client.Get("https://api.example.com/orders/1")
```

出站请求自动携带 `traceparent`，对端按 W3C 标准即可关联。

## 日志联动

```go
logger.Info("处理订单", tracex.LogFields(ctx))
// {"trace_id":"...","span_id":"..."}
```

## 全局集成

使用 OTel 全局 API 的库（如第三方 SDK）可通过 `SetGlobal: true`
一键注册：TracerProvider 与传播器变为全局默认，库代码直接调用
`otel.Tracer(...)` / `otel.GetTextMapPropagator()` 即可。
