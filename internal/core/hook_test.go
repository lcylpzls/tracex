package core

import (
	"context"
	"errors"
	testx "github.com/lcylpzls/testx"
	"io"
	"testing"
	"time"

	"github.com/lcylpzls/logx"
)

// TestRecordError 覆盖错误记录分支。
func TestRecordError(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	ctx, span := m.Start(context.Background(), "op")
	RecordError(ctx, errors.New("boom"))
	RecordError(ctx, nil) // nil 忽略
	RecordError(context.Background(), errors.New("ignored"))
	span.End()
	_ = m.Shutdown(context.Background())
	spans := m.Spans()
	if len(spans) != 1 {
		t.Fatalf("应收集 1 条 span：%+v", spans)
	}
	s := spans[0]
	testx.RequireEqual(t, s.StatusCode, "Error")

	found := false
	for _, ev := range s.Events {
		if ev.Name == "exception" {
			found = true
		}
	}
	testx.RequireTrue(t, found)

}

// TestLogHook 覆盖日志写入 span 事件。
func TestLogHook(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	logger, err := logx.NewBuilder().EnableWriter(io.Discard, logx.InfoLevel).Build()
	testx.RequireNoError(t, err)

	hl, ok := logger.(logx.HookedLogger)
	testx.RequireTrue(t, ok)

	hl.AddHook(NewLogHook())

	ctx, span := m.Start(context.Background(), "op")
	logger.WithContext(ctx).Info("hello", logx.Fields())
	logger.WithContext(context.Background()).Info("noop", logx.Fields())
	time.Sleep(100 * time.Millisecond) // 等待异步钩子触发
	span.End()
	_ = m.Shutdown(context.Background())

	spans := m.Spans()
	if len(spans) != 1 {
		t.Fatalf("应收集 1 条 span：%+v", spans)
	}
	var logEvents []SpanEvent
	for _, ev := range spans[0].Events {
		if ev.Name == "log" {
			logEvents = append(logEvents, ev)
		}
	}
	if len(logEvents) != 1 {
		t.Fatalf("应记录 1 条日志事件：%+v", spans[0].Events)
	}
	if logEvents[0].Attributes["log.message"] != "hello" ||
		logEvents[0].Attributes["log.level"] == "" {
		t.Fatalf("日志事件内容不符：%+v", logEvents[0])
	}
}
