# 更新日志

## [v0.2.0] - 2026-08-10

### 新增

- Baggage 便捷 API：`WithBaggage` / `BaggageValue`；
- 出站 HTTP 注入：`Manager.RoundTripper` 自动携带 traceparent、
  记录客户端 span，并按状态码/错误标记结果；
- `AddSpanEvent` span 事件记录；
- Span 快照支持事件；`FuzzBaggage` 模糊目标接入 CI。

### 质量

- 语句覆盖率 100%；race / vet / staticcheck / fuzz 全绿。

## [v0.1.0] - 2026-08-10

### 新增

- TracerProvider 管理：stdout / OTLP/HTTP / 内存导出器；
- 标准 net/http 链路追踪中间件（traceparent 透传、状态码与
  5xx 错误标记）；
- `LogFields(ctx)` logx 字段联动；
- 内存导出器与 Span 快照（测试/调试）；
- errx 错误码全集；
- 三平台 + Linux 多发行版 CI/Release。

### 质量

- 语句覆盖率 100%；vet / staticcheck / race 全绿。
