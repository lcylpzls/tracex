package filex

import (
	"context"
	"errors"
	testx "github.com/lcylpzls/testx"
	"testing"

	"github.com/lcylpzls/filex"
	"github.com/lcylpzls/tracex"
)

// TestHook 覆盖 filex 追踪钩子成功与失败路径。
func TestHook(t *testing.T) {
	m, err := tracex.New(tracex.Config{ServiceName: "svc", Exporter: tracex.ExporterMemory})
	testx.RequireNoError(t, err)

	h := NewHook(m)
	ctx, end := h.Start(context.Background(), "filex.put",
		filex.TraceAttr{Key: "filex.operation", Value: "put"})
	end(nil)
	ctx2, end2 := h.Start(context.Background(), "filex.get",
		filex.TraceAttr{Key: "filex.operation", Value: "get"})
	end2(errors.New("对象不存在"))
	_ = ctx
	_ = ctx2
	_ = m.Shutdown(context.Background())

	spans := m.Spans()
	if len(spans) != 2 {
		t.Fatalf("应收集 2 条 span：%+v", spans)
	}
	ok, fail := 0, 0
	for _, s := range spans {
		if s.Attributes["filex.operation"] == "" {
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
