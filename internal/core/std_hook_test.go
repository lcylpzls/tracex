package core

import (
	"context"
	"errors"
	"testing"
)

func TestNewHookStart(t *testing.T) {
	m, err := New(Config{ServiceName: "smoke", Exporter: ExporterMemory})
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}
	h := NewHook(m)
	if h == nil {
		t.Fatal("NewHook 返回 nil")
	}
	ctx, end := h.Start(context.Background(), "op", TraceAttr{Key: "k", Value: "v"})
	end(nil)
	end(errors.New("模拟失败"))
	_ = ctx
}
