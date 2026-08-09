package tracex

import (
	"context"
	"testing"
)

// FuzzBaggage 保证任意键值写入/读取不 panic。
func FuzzBaggage(f *testing.F) {
	f.Add("key", "value")
	f.Add("", "")
	f.Add("bad key!", "x")
	f.Fuzz(func(t *testing.T, key, value string) {
		ctx := WithBaggage(context.Background(), key, value)
		_ = BaggageValue(ctx, key)
	})
}
