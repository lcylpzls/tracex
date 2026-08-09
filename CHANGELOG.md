# 更新日志

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
