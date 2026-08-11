package tracex_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lcylpzls/tracex"
)

// TestPublicAPI 黑盒冒烟测试：覆盖根包全部转发函数、类型别名与常量。
func TestPublicAPI(t *testing.T) {
	if tracex.Version != "v1.3.0" {
		t.Fatalf("Version 不符：%s", tracex.Version)
	}

	m, err := tracex.New(tracex.Config{ServiceName: "smoke", Exporter: tracex.ExporterMemory})
	if err != nil || m == nil {
		t.Fatalf("New 失败：%v", err)
	}

	hook := tracex.NewHook(m)
	if hook == nil {
		t.Fatal("NewHook 返回 nil")
	}
	ctx, end := hook.Start(context.Background(), "smoke_op", tracex.TraceAttr{Key: "k", Value: "v"})
	end(nil)
	end(errors.New("失败"))
	_ = ctx

	lh := tracex.NewLogHook()
	if lh == nil {
		t.Fatal("NewLogHook 返回 nil")
	}

	ctx2 := tracex.WithBaggage(context.Background(), "key", "value")
	if tracex.BaggageValue(ctx2, "key") != "value" {
		t.Fatal("BaggageValue 不一致")
	}

	me := tracex.NewMemoryExporter()
	if me == nil {
		t.Fatal("NewMemoryExporter 返回 nil")
	}
	_ = tracex.LogFields(ctx2)
	tracex.AddSpanEvent(ctx2, "event")
	tracex.RecordError(ctx2, errors.New("记录错误"))

	_ = tracex.ExporterStdout
	_ = tracex.ExporterOTLPHTTP
	_ = tracex.CodeInvalidConfig
	_ = tracex.CodeExporterFailed
	_ = tracex.CodeShutdownFailed
	var _ tracex.Config
	var _ tracex.SpanSnapshot
	var _ tracex.SpanEvent
	var _ tracex.TraceHook
	var _ tracex.TraceAttr
	var _ tracex.LogHook
	var _ tracex.Manager
	var _ tracex.MemoryExporter
}
