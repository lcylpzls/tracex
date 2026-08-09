# 更新日志

## [v0.4.0] - 2026-08-10

### 新增

- OTLP/HTTP 超时可配置（`OTLPTimeout`，默认 10s）；
- 采样器可插拔（`Config.Sampler`，nil 使用采样率）；
- 并发安全回归测试（32 并发入站/出站，race 检测）；
- CI 增加 go.mod tidy 漂移检查。

### 质量

- 语句覆盖率保持 100%；race / vet / staticcheck / fuzz 全绿。

## [v0.3.0] - 2026-08-10

### 新增

- 路由级 span 命名：`Config.RouteNamer` 支持框架适配注入
  路由模板（如 `/users/{id}`）；
- 慢请求标记：`Config.SlowThreshold` 超时后记录 `slow` 事件与
  耗时属性。

### 质量

- 语句覆盖率保持 100%；race / vet / staticcheck 全绿。

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
