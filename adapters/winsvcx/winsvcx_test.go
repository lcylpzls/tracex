//go:build windows

package winsvcx

import (
	"context"
	"errors"
	"testing"

	"github.com/lcylpzls/tracex"
	"github.com/lcylpzls/winsvcx/lib/service"
)

// TestHook 覆盖 winsvcx 追踪钩子成功与失败路径。
func TestHook(t *testing.T) {
	m, err := tracex.New(tracex.Config{ServiceName: "svc", Exporter: tracex.ExporterMemory})
	if err != nil {
		t.Fatal(err)
	}
	h := NewHook(m)
	ctx, end := h.Start(context.Background(), "winsvcx.service.execute",
		service.TraceAttr{Key: "winsvcx.service_name", Value: "MyService"})
	end(nil)
	ctx2, end2 := h.Start(context.Background(), "winsvcx.service.execute",
		service.TraceAttr{Key: "winsvcx.service_name", Value: "MyService"})
	end2(errors.New("服务退出码：1"))
	_ = ctx
	_ = ctx2
	_ = m.Shutdown(context.Background())

	spans := m.Spans()
	if len(spans) != 2 {
		t.Fatalf("应收集 2 条 span：%+v", spans)
	}
	ok, fail := 0, 0
	for _, s := range spans {
		if s.Attributes["winsvcx.service_name"] != "MyService" {
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
