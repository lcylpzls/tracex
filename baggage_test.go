package tracex

import (
	"context"
	"testing"
)

// TestBaggage 覆盖 bag 成员写入与读取。
func TestBaggage(t *testing.T) {
	ctx := WithBaggage(context.Background(), "user_id", "42")
	if got := BaggageValue(ctx, "user_id"); got != "42" {
		t.Fatalf("bag 值不符：%q", got)
	}
	if got := BaggageValue(ctx, "missing"); got != "" {
		t.Fatalf("缺失键应返回空串：%q", got)
	}
	// 非法键被忽略，原上下文不受影响。
	orig := WithBaggage(context.Background(), "ok", "1")
	after := WithBaggage(orig, "bad key!", "x")
	if BaggageValue(after, "ok") != "1" {
		t.Fatal("非法键不应破坏已有 bag")
	}
	if BaggageValue(after, "bad key!") != "" {
		t.Fatal("非法键不应写入")
	}
}
