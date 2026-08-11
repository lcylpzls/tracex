package core

import "github.com/lcylpzls/tracex/contract"

// TraceAttr 是链路追踪属性键值对（契约定义见 tracex/contract）。
type TraceAttr = contract.TraceAttr

// TraceHook 是家族库通用的链路追踪钩子契约（定义见 tracex/contract）：
// 各库只依赖该接口，不绑定任何追踪后端；tracex 实现负责接入
// OpenTelemetry 等具体后端。
type TraceHook = contract.TraceHook
