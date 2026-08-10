# 更新日志

## [v1.0.3] - 2026-08-10

### 变更

- 家族正式基线锁定：依赖统一指向 v1 基线已发布版本（errx v1.5.5 / logx v1.3.2 / testx v1.4.3 / validx v1.2.4 / cryptox v1.0.2 / confx v1.0.2 / webx v1.5.4 等），此后家族依赖不再前进。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.0.2] - 2026-08-10

### 变更

- 家族依赖最终对齐到 v1 正式版基线（errx v1.5.4 / logx v1.3.1 / testx v1.4.2 / validx v1.2.3 / confx v1.0.1 / cryptox v1.0.1 等），无 API 变更。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.0.1] - 2026-08-10

### 变更

- 家族依赖统一对齐到最新基线（errx v1.5.4 / logx v1.3.0 / testx v1.4.1 / validx v1.2.2 等），无 API 变更。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.0.0] - 2026-08-10

### 发布

- 家族正式版 v1.0.0：当前 API 与行为作为 v1 基线；按家族规则，v1.*.* 内允许破坏性修改，不承诺向后兼容。

### 质量

- 覆盖率、race、vet、staticcheck、fuzz、govulncheck 全绿；CI/Release 自动化发布。

## [v0.17.1] - 2026-08-10

### 变更

- go 指令与 CI/Release 工作流统一为 Go 1.26.5；
- README Go 版本徽章同步更新。

## [v0.17.0] - 2026-08-10

### 变更

- 家族测试底座接入：根包与全部 adapters 测试改用语义等价的
  testx 断言（含 Require* 致命断言）；
- 测试依赖新增 `testx v1.2.0`，errx 同步升级 v1.4.0；
- adapters 子模块同步发布 `adapters/v0.17.0`。

### 质量

- 根包与全部适配包语句覆盖率 100%；race / vet / staticcheck 全绿。

## [v0.16.0] - 2026-08-10

### 变更

- adapters 子模块升级依赖：webx → `webx/v2 v2.0.1`（v2 主版本
  模块路径）、logx v1.1.0、errx v1.4.0；
- adapters 子模块发布机制补齐：新增 `adapters/vX.Y.Z` tag 支持，
  子模块可被 `go get github.com/lcylpzls/tracex/adapters` 消费。

## [v0.15.0] - 2026-08-10

### 新增

- adapters 新增 filex（对象存储）追踪钩子适配，配合 filex v0.20.0
  （Put/Get/Head/Delete/List/Copy/Move 埋点）。

### 质量

- 全部 9 个适配器语句覆盖率 100%；race / vet / staticcheck 全绿。

## [v0.14.0] - 2026-08-10

### 家族插拔（收官）

- adapters 新增 authx（刷新令牌操作）与 winsvcx（服务生命周期）
  追踪钩子适配；配合家族同步发版：
  - authx v1.0.4：refresh issue/validate/consume/rotate 埋点；
  - winsvcx v0.13.0：RunWithHook 服务会话埋点；
  - webx v1.2.6：请求 ID 头名文档统一为 X-Request-ID；
- 新增家族可观测性规范文档（接入矩阵 / 不接入清单与理由 /
  全栈接入示例）。

### 质量

- 全部 8 个适配器语句覆盖率 100%；race / vet / staticcheck 全绿。

## [v0.13.0] - 2026-08-10

### 规范落地（方案 B）

- 配合 webx v1.2.5：全局中间件链覆盖 404/405 兜底请求，
  webx 适配器追踪无盲区（补 404 span 断言）；
- 明确规范：webx 不内置分布式追踪，链路追踪统一由 tracex 基座 +
  adapters 接入。

### 质量

- 全部 6 个适配器语句覆盖率 100%；race / vet / staticcheck 全绿。

## [v0.12.0] - 2026-08-10

### 家族插拔（续）

- adapters 新增 resiliencex / updatex 追踪钩子适配；配合家族同步
  发版：
  - resiliencex v1.0.3：熔断执行埋点（ExecuteContext）；
  - updatex v0.5.0：Check/Apply 埋点；
- 适配器集成文档补全（现覆盖 webx/dbx/jobx/cachex/resiliencex/
  updatex/httpx）。

### 质量

- 全部 6 个适配器语句覆盖率 100%；race / vet / staticcheck 全绿。

## [v0.11.0] - 2026-08-10

### 家族插拔

- adapters 新增 dbx / jobx / cachex 追踪钩子适配（零核心依赖，
  保持 tracex 核心薄）；配合家族同步发版：
  - dbx v0.2.4：Exec/Query/QueryRow 埋点；
  - jobx v1.0.4：任务执行埋点；
  - cachex v1.0.3：回源加载埋点；
  - httpx v1.0.4：`WithRoundTripperWrapper` 传输层插拔；
- 适配器集成文档（webx/dbx/jobx/cachex/httpx 用法）。

### 质量

- 全部适配器语句覆盖率 100%；race / vet / staticcheck 全绿。

## [v0.10.0] - 2026-08-10

### 终轮

- OTLP Collector 部署文档（collector 配置/Docker/tracex 接入/
  生产建议）；
- 最终审计全量通过（三平台 + 多发行版 + 交叉构建 + bench + fuzz +
  tidy + govulncheck）。

> **自评：tracex 已达到 v1 候选标准（第二轮）**，v1 是否发布由用户决定。

## [v0.9.0] - 2026-08-10

### 新增

- webx 框架适配（独立 adapters 模块）：全局中间件自动提取/创建
  链路、路由级命名、状态码与慢请求标记；
- adapters 端到端集成测试（webx 服务 + 内存导出器，100% 覆盖）
  接入三平台 CI 与 Release。

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
