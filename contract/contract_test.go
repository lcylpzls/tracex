package contract

import (
	"context"
	"testing"
)

// hookFunc 是 TraceHook 的函数式实现（测试辅助）。
type hookFunc func(context.Context, string, ...TraceAttr) (context.Context, func(error))

func (f hookFunc) Start(ctx context.Context, name string, attrs ...TraceAttr) (context.Context, func(error)) {
	return f(ctx, name, attrs...)
}

// TestContract 覆盖契约类型与 Start 语义。
func TestContract(t *testing.T) {
	var h TraceHook = hookFunc(func(ctx context.Context, _ string, _ ...TraceAttr) (context.Context, func(error)) {
		return ctx, func(error) {}
	})
	ctx, end := h.Start(context.Background(), "op", TraceAttr{Key: "k", Value: "v"})
	if ctx == nil {
		t.Fatal("ctx 不应为 nil")
	}
	end(nil)
}
