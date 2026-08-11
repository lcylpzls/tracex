# tracex 架构详解

> 版本：v1.3.0（已实现，随版本演进）

## 1. 模块结构

```
tracex/
├── tracex.go      # Config / Manager / 导出器构造
├── middleware.go  # 标准 net/http 链路追踪中间件
├── logfields.go   # logx 日志字段联动
├── memory.go      # 内存导出器（测试/调试）
└── errors.go      # errx 错误码
```

## 2. Manager

`Manager` 持有 TracerProvider、导出器、组合传播器
（TraceContext + Baggage）与 Tracer，生命周期：

- `New(Config)`：校验配置 → 构造导出器 → 组装 Provider；
- `Shutdown(ctx)`：刷新并关闭（幂等）。

## 3. 中间件

`Manager.Middleware(next)`：

1. 从请求头提取 traceparent（无则创建新链路）；
2. 以 URL path 为 span 名，记录 method/path/scheme/host/UA；
3. 包装 ResponseWriter 捕获状态码；
4. 5xx 标记 `Error` 状态；
5. End 后由批量处理器异步导出。

## 4. 导出器

- `stdout`：默认，开发调试；
- `memory`：测试/调试，`Manager.Spans()` 读取快照；
- `otlp-http`：生产采集，支持端点/明文/自定义头。

## 5. 日志联动

`LogFields(ctx)` 提取 SpanContext 的 trace_id / span_id，
返回 logx 字段组；无有效链路时为空字段组。
