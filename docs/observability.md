# 家族可观测性规范

## 总则

1. **唯一追踪基座**：分布式追踪只由 tracex 实现（OTel），其余
   基座禁止直接依赖 OTel；
2. **零依赖插拔**：需要追踪的基座只暴露最小 `TraceHook` 接口
   （默认 no-op）或传输包装点，由 tracex/adapters 统一接入；
3. **无内置追踪**：任何基座不内置 span 逻辑；`X-Request-ID`
   一类请求 ID 属于请求标识，与链路追踪无关。

## 接入矩阵

| 基座 | 接入点 | 适配器 | 版本 |
| --- | --- | --- | --- |
| webx | 全局中间件 | tracex/adapters/webx | v1.2.5+ |
| httpx | WithRoundTripperWrapper | m.RoundTripper | v1.0.4+ |
| dbx | WithTraceHook | tracex/adapters/dbx | v0.2.4+ |
| jobx | WithTraceHook | tracex/adapters/jobx | v1.0.4+ |
| cachex | WithTraceHook | tracex/adapters/cachex | v1.0.3+ |
| resiliencex | WithTraceHook + ExecuteContext | tracex/adapters/resiliencex | v1.0.3+ |
| updatex | Config.TraceHook | tracex/adapters/updatex | v0.5.0+ |
| authx | token.WithTraceHook | tracex/adapters/authx | v1.0.4+ |
| winsvcx | service.RunWithHook | tracex/adapters/winsvcx | v0.13.0+ |
| filex | Config.TraceHook | tracex/adapters/filex | v0.20.0+ |

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
m, _ := tracex.New(tracex.Config{ServiceName: "order", Exporter: tracex.ExporterOTLPHTTP, ...})

// webx 入站
s.UseGlobalMiddleware(txwebx.Middleware(m))
// httpx 出站
client, _ := httpx.New(httpx.WithRoundTripperWrapper(m.RoundTripper))
// dbx / jobx / cachex / resiliencex / updatex / authx / winsvcx
dbx.Open(ctx, "mysql", dsn, dbx.WithTraceHook(txdbx.NewHook(m)))
jobx.NewDispatcher(jobx.WithTraceHook(txjobx.NewHook(m)))
cachex.New(cachex.WithTraceHook(txcachex.NewHook(m)))
resiliencex.NewCircuitBreaker(resiliencex.WithTraceHook(txresiliencex.NewHook(m)))
updatex.New(updatex.Config{TraceHook: txupdatex.NewHook(m), ...})
token.IssueRefreshToken(ctx, store, ttl, token.WithTraceHook(txauthx.NewHook(m)))
service.RunWithHook(name, txwinsvcx.NewHook(m))
```
