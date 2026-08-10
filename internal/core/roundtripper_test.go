package core

import (
	"context"
	testx "github.com/lcylpzls/testx"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// TestRoundTripper 覆盖出站注入与状态标记。
func TestRoundTripper(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	okSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer okSrv.Close()
	client := &http.Client{Transport: m.RoundTripper(http.DefaultTransport)}
	if _, err := client.Get(okSrv.URL + "/ok"); err != nil {
		t.Fatalf("请求失败：%v", err)
	}

	errSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer errSrv.Close()
	if _, err := client.Get(errSrv.URL + "/err"); err != nil {
		t.Fatalf("请求失败：%v", err)
	}

	// 默认传输（nil 分支）。
	client2 := &http.Client{Transport: m.RoundTripper(nil)}
	if _, err := client2.Get(okSrv.URL + "/nil"); err != nil {
		t.Fatalf("默认传输请求失败：%v", err)
	}

	// 网络错误。
	dead := httptest.NewServer(nil)
	deadURL := dead.URL
	dead.Close()
	if _, err := client.Get(deadURL + "/dead"); err == nil {
		t.Fatal("应发生连接错误")
	}
	_ = m.Shutdown(context.Background())

	spans := m.Spans()
	var okSpan, errSpan, nilSpan, deadSpan SpanSnapshot
	for _, s := range spans {
		switch {
		case strings.HasSuffix(s.Name, "/ok"):
			okSpan = s
		case strings.HasSuffix(s.Name, "/err"):
			errSpan = s
		case strings.HasSuffix(s.Name, "/nil"):
			nilSpan = s
		case strings.HasSuffix(s.Name, "/dead"):
			deadSpan = s
		}
	}
	if okSpan.StatusCode != "Ok" || okSpan.Attributes["http.response.status_code"] != "200" {
		t.Fatalf("OK span 不符：%+v", okSpan)
	}
	if errSpan.StatusCode != "Error" || errSpan.Attributes["http.response.status_code"] != "502" {
		t.Fatalf("5xx span 不符：%+v", errSpan)
	}
	testx.RequireNotEqual(t, nilSpan.Name, "")

	if deadSpan.StatusCode != "Error" || len(deadSpan.Events) == 0 {
		t.Fatalf("网络错误 span 不符：%+v", deadSpan)
	}
}

// TestAddSpanEvent 覆盖 span 事件记录。
func TestAddSpanEvent(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	ctx, span := m.Start(context.Background(), "op")
	AddSpanEvent(ctx, "cache.hit", attribute.String("key", "x"))
	AddSpanEvent(context.Background(), "ignored") // 无活动 span 应忽略
	span.End()
	_ = m.Shutdown(context.Background())
	spans := m.Spans()
	if len(spans) != 1 || len(spans[0].Events) != 1 {
		t.Fatalf("事件快照不符：%+v", spans)
	}
	ev := spans[0].Events[0]
	if ev.Name != "cache.hit" || ev.Attributes["key"] != "x" {
		t.Fatalf("事件内容不符：%+v", ev)
	}
}
