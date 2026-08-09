package updatex

import (
	"context"
	"errors"
	"testing"

	"github.com/lcylpzls/tracex"
	"github.com/lcylpzls/updatex"
)

// TestHook 覆盖 updatex 追踪钩子成功与失败路径。
func TestHook(t *testing.T) {
	m, err := tracex.New(tracex.Config{ServiceName: "svc", Exporter: tracex.ExporterMemory})
	if err != nil {
		t.Fatal(err)
	}
	h := NewHook(m)
	ctx, end := h.Start(context.Background(), "updatex.check",
		updatex.TraceAttr{Key: "updatex.current_version", Value: "1.0.0"})
	end(nil)
	ctx2, end2 := h.Start(context.Background(), "updatex.apply",
		updatex.TraceAttr{Key: "updatex.current_version", Value: "1.0.0"})
	end2(errors.New("下载失败"))
	_ = ctx
	_ = ctx2
	_ = m.Shutdown(context.Background())

	spans := m.Spans()
	if len(spans) != 2 {
		t.Fatalf("应收集 2 条 span：%+v", spans)
	}
	ok, fail := 0, 0
	for _, s := range spans {
		if s.Attributes["updatex.current_version"] != "1.0.0" {
			t.Fatalf("属性不符：%+v", s.Attributes)
		}
		if s.StatusCode == "Error" {
			fail++
		} else {
			ok++
		}
	}
	if ok != 1 || fail != 1 {
		t.Fatalf("状态计数不符：ok=%d fail=%d", ok, fail)
	}
}
