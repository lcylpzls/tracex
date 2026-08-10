package core

import "context"

// TraceAttr 是链路追踪属性键值对，家族库 TraceHook 契约的通用结构。
type TraceAttr struct {
	Key   string
	Value string
}

// TraceHook 是家族库通用的链路追踪钩子契约：
// 各库只依赖该接口，不绑定任何追踪后端；tracex 适配器负责接入
// OpenTelemetry 等具体实现。
type TraceHook interface {
	// Start 在操作开始前调用：返回携带链路上下文的 ctx 与结束回调
	// （结束回调入参为操作结果错误，nil 表示成功）。
	Start(ctx context.Context, name string, attrs ...TraceAttr) (context.Context, func(error))
}
