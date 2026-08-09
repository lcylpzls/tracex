# tracex 版本路线

> 目标：v0.1.0 起每版完成即全自动 CI + Release，全部通过后进入下一版；
> 继续自我检查自我提升推进版本号，直到自评达到 v1 候选后停下，
> v1 是否发布由用户决定。

## v0.1.0 — 核心闭环

- TracerProvider 管理（stdout / OTLP/HTTP / 内存导出器）；
- 标准 net/http 链路追踪中间件（traceparent 透传、状态码与 5xx 错误标记）；
- `LogFields(ctx)` logx 字段联动；
- errx 错误码全集与 100% 语句覆盖；
- 三平台 + Linux 多发行版 CI，Release 全自动。

## v0.2.0 — 家族框架适配

- 出站 HTTP 注入（RoundTripper）；
- Baggage 注入/读取便捷 API；
- span 事件记录与快照。

> 状态：**已发布**（v0.2.0，2026-08-10）。

## v0.3.0 — 可观测增强

- 路由级 span 命名（RouteNamer）；
- 慢请求标记（SlowThreshold + slow 事件）。

> 状态：**已发布**（v0.3.0，2026-08-10）。

## v0.4.0 — 健壮性打磨

- OTLP 重连与超时配置；采样器可插拔；
- 并发与竞态终检（race + fuzz）；
- 依赖收敛审计。

> 状态：**已发布**（v0.4.0，2026-08-10）。

## v0.5.0 — 发布前终审

- 文档定稿、LICENSE、govulncheck、三架构构建；
- 收口 roadmap，进入自主打磨。

> 状态：**已发布**（v0.5.0，2026-08-10）。交叉构建矩阵
> （3 平台 × 2 架构）接入 CI/Release。

## v0.6.0+ — 自主打磨

- 持续自我审查与真实缺陷修复，直到 v1 候选。

### v0.6.0（已发布）

- 可选注册 OTel 全局组件（SetGlobal）。

### v0.7.0（已发布，v1 候选）

- 集成示例文档；最终审计全绿。

### v0.8.0（已发布）

- RecordError / LogHook / 可运行示例 / 基准测试。

### v0.9.0（已发布）

- webx 框架适配（adapters 子模块，端到端集成测试）。

### v0.10.0（已发布，v1 候选）

- OTLP 部署文档；最终审计全绿。

### v0.11.0（已发布）

- 家族插拔：dbx/jobx/cachex 钩子适配 + httpx 传输层包装；
  各库同步发版。

### v0.12.0（已发布）

- 家族插拔续：resiliencex 熔断执行 + updatex 自更新适配。

> 自评：v1 候选达成（第二轮），停止自动推进；v1 是否发布由用户决定。

> 自评：v1 候选达成，停止自动推进；v1 是否发布由用户决定。

## 质量门禁（每版）

```powershell
go test -count=1 ./...
go test -count=1 -coverprofile=coverage.out ./...   # 100%
go test -race -count=1 ./...
go vet ./... && staticcheck ./...
```

CI：ubuntu/windows/macos 三平台 + Linux 多发行版容器矩阵 +
govulncheck；Release（tag 触发，全绿后发布）。
