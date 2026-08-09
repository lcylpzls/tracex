# tracex 设计定版

> 版本：v0.1.0

## 1. 定位

tracex 是 OpenTelemetry 链路追踪基础库，解决服务接入分布式追踪时的
重复劳动：TracerProvider 初始化、HTTP 入站/出站上下文透传、日志关联。

**边界**：

- 只做追踪（trace），不做指标/日志采集；
- 不绑定具体框架：HTTP 中间件基于标准 net/http，框架适配放子包；
- 生产采集使用 OTLP/HTTP，调试使用 stdout/内存导出器。

## 2. 依赖

- `go.opentelemetry.io/otel`（API）与 `otel/sdk`（实现）——必需；
- OTLP/HTTP、stdout 导出器——必需；
- `errx` / `logx`——家族错误与日志。

## 3. 数据流

```
入站请求 ──► Middleware
              ├─ Extract(traceparent) ──► 关联父链路
              ├─ Start(span) ──► 记录 method/path/status
              ├─ 5xx ──► SetStatus(Error)
              └─ End ──► BatchProcessor ──► Exporter(stdout/OTLP/内存)

业务代码 ──► LogFields(ctx) ──► logx 输出 trace_id/span_id
出站请求 ──► Inject(ctx, carrier) ──► 携带 traceparent
```

## 4. 错误码

| 错误码 | 分类 |
| --- | --- |
| tracex_invalid_config | invalid |
| tracex_exporter_failed | unavailable |
| tracex_shutdown_failed | unavailable |

## 5. 质量目标

- 语句覆盖率 100%；race / vet / staticcheck / fuzz；
- 三平台 + Linux 多发行版 CI。
