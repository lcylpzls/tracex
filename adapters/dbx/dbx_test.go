package dbx

import (
	"context"
	"errors"
	testx "github.com/lcylpzls/testx"
	"testing"

	"github.com/lcylpzls/dbx"
	"github.com/lcylpzls/tracex"
)

// TestHook 覆盖 dbx 追踪钩子成功与失败路径。
func TestHook(t *testing.T) {
	m, err := tracex.New(tracex.Config{ServiceName: "svc", Exporter: tracex.ExporterMemory})
	testx.RequireNoError(t, err)

	h := NewHook(m)

	ctx, end := h.Start(context.Background(), "dbx.exec",
		dbx.TraceAttr{Key: "db.operation", Value: "exec"})
	end(nil)
	ctx2, end2 := h.Start(context.Background(), "dbx.query",
		dbx.TraceAttr{Key: "db.operation", Value: "query"})
	end2(errors.New("查询失败"))
	_ = ctx
	_ = ctx2
	_ = m.Shutdown(context.Background())

	spans := m.Spans()
	if len(spans) != 2 {
		t.Fatalf("应收集 2 条 span：%+v", spans)
	}
	var okSpan, errSpan tracex.SpanSnapshot
	for _, s := range spans {
		if s.Attributes["db.operation"] == "exec" {
			okSpan = s
		} else {
			errSpan = s
		}
	}
	if okSpan.StatusCode != "Unset" || errSpan.StatusCode != "Error" {
		t.Fatalf("状态不符：ok=%s err=%s", okSpan.StatusCode, errSpan.StatusCode)
	}
	exception := false
	for _, ev := range errSpan.Events {
		if ev.Name == "exception" {
			exception = true
		}
	}
	testx.RequireTrue(t, exception)

}
