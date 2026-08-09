package tracex

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

// WithBaggage 返回携带 bag 成员的上下文；非法键值被忽略。
func WithBaggage(ctx context.Context, key, value string) context.Context {
	member, err := baggage.NewMember(key, value)
	if err != nil {
		return ctx
	}
	bag, _ := baggage.New(member)
	return baggage.ContextWithBaggage(ctx, bag)
}

// BaggageValue 读取上下文中的 bag 成员值；不存在返回空串。
func BaggageValue(ctx context.Context, key string) string {
	return baggage.FromContext(ctx).Member(key).Value()
}
