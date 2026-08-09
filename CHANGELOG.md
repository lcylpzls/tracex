# 更新日志

## [v0.8.0] - 2026-08-10

### 新增

- `RecordError` 错误记录（异常事件 + Error 状态）；
- `LogHook` logx 日志钩子：日志写入当前 span 的 log 事件，
  实现日志与链路事件关联；
- 可运行示例 `examples/demo`（stdout + 中间件 + 出站注入 + 钩子）；
- 基准测试（中间件 ~1.8µs/op、出站注入含 HTTP 往返）接入 CI bench。

### 质量

- 语句覆盖率保持 100%；race / vet / staticcheck / fuzz 全绿。

## [v0.7.0] - 2026-08-10

### 终轮

- 集成示例文档（服务端/客户端/日志/全局集成）；
- 最终审计：三平台 × 双架构、race、shuffle、fuzz、vet、
  staticcheck、依赖校验全绿。

> **自评：tracex 已达到 v1 候选标准**，v1 是否发布由用户决定。

## [v0.6.0] - 2026-08-10

### 打磨

- 可选注册 OTel 全局组件：`Config.SetGlobal` 将 TracerProvider 与
  传播器注册为全局默认，方便使用 otel 全局 API 的库无缝接入。

### 质量

- 语句覆盖率保持 100%；race / vet / staticcheck 全绿。

## [v0.5.0] - 2026-08-10

### 发布前终审

- CI 增加三平台 × 双架构（linux/windows/darwin × amd64/arm64）
  交叉构建；
- Release 增加 6 组交叉构建验证；
- 文档定稿、依赖收敛审计（go mod tidy/verify）。

> roadmap 至此完成，后续版本为自主打磨。

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
