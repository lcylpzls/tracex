package tracex

import (
	"context"
	"github.com/lcylpzls/testx"
	"testing"
)

// TestBaggage 覆盖 bag 成员写入与读取。
func TestBaggage(t *testing.T) {
	ctx := WithBaggage(context.Background(), "user_id", "42")
	testx.RequireEqual(t, BaggageValue(ctx, "user_id"), "42")
	testx.RequireEqual(t, BaggageValue(ctx, "missing"), "")
	// 非法键被忽略，原上下文不受影响。
	orig := WithBaggage(context.Background(), "ok", "1")
	after := WithBaggage(orig, "bad key!", "x")
	testx.RequireEqual(t, BaggageValue(after, "ok"), "1")
	testx.RequireEqual(t, BaggageValue(after, "bad key!"), "")
}
