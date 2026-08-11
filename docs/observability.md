# 家族可观测性规范

## 总则

1. **唯一追踪基座**：分布式追踪只由 tracex 实现（OTel），其余
   基座禁止直接依赖 OTel；
2. **零依赖插拔**：需要追踪的基座只暴露最小 `TraceHook` 接口
   （默认 no-op）或传输包装点，由调用方用
   `tracex.NewHook(m)` 统一接入；
3. **无内置追踪**：任何基座不内置 span 逻辑；`X-Request-ID`
   一类请求 ID 属于请求标识，与链路追踪无关。

## 接入矩阵

| 基座 | 接入点 | 接入方式 |
| --- | --- | --- |
| webx | 全局中间件 | `s.UseGlobalMiddleware(m.Middleware)` |
| httpx | WithRoundTripperWrapper | `m.RoundTripper` |
| dbx | WithTraceHook | `tracex.NewHook(m)` |
| jobx | WithTraceHook | `tracex.NewHook(m)` |
| cachex | WithTraceHook | `tracex.NewHook(m)` |
| resiliencex | WithTraceHook | `tracex.NewHook(m)` |
| updatex | Config.TraceHook | `tracex.NewHook(m)` |
| authx | token.WithTraceHook | `tracex.NewHook(m)` |
| winsvcx | service.RunWithHook | `tracex.NewHook(m)` |
| filex | Config.TraceHook | `tracex.NewHook(m)` |

## 不接入清单与理由

| 基座 | 理由 |
| --- | --- |
| confx | 纯同步配置解析，无跨进程/IO 边界，埋点零收益 |
| validx | 纯同步校验，微秒级，无边界 |
| idgenx | 纯同步 ID 生成，无边界 |
| errx | 错误类型定义，非执行单元 |
| logx | 日志底座；与 tracex 反向联动（LogFields / LogHook 已覆盖） |
| metricsx | 指标输出适配，非追踪对象 |

## 接入示例（全栈链路）

```go
m, _ := tracex.New(tracex.Config{
	ServiceName: "order",
	Exporter:    tracex.ExporterOTLPHTTP,
	// ...
})
hook := tracex.NewHook(m)

// webx 入站（标准中间件）
s.UseGlobalMiddleware(m.Middleware)
// httpx 出站
client, _ := httpx.New(httpx.WithRoundTripperWrapper(m.RoundTripper))
// dbx / jobx / cachex / resiliencex / updatex / authx / winsvcx / filex
dbx.Open(ctx, "mysql", dsn, dbx.WithTraceHook(hook))
jobx.NewDispatcher(jobx.WithTraceHook(hook))
cachex.New(cachex.WithTraceHook(hook))
resiliencex.NewCircuitBreaker(resiliencex.WithTraceHook(hook))
updatex.New(updatex.Config{TraceHook: hook})
token.IssueRefreshToken(ctx, store, ttl, token.WithTraceHook(hook))
service.RunWithHook(name, hook) // winsvcx（Windows 平台）
filex.New(filex.Config{TraceHook: hook})
```
