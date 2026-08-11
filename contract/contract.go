// Package contract 定义家族统一的链路追踪契约（零第三方依赖）。
// 家族库只依赖本包，不绑定任何追踪后端；tracex 等实现通过实现
// TraceHook 接入 OpenTelemetry 或其它后端。
package contract

import "context"

// TraceAttr 是链路追踪属性键值对。
type TraceAttr struct {
	Key   string
	Value string
}

// TraceHook 是家族库通用的链路追踪钩子契约：
// Start 在操作开始前调用，返回携带链路上下文的 ctx 与结束回调
// （结束回调入参为操作结果错误，nil 表示成功）。
type TraceHook interface {
	Start(ctx context.Context, name string, attrs ...TraceAttr) (context.Context, func(error))
}
