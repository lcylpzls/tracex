package tracex

import (
	"bytes"
	"context"
	"fmt"
	testx "github.com/lcylpzls/testx"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lcylpzls/errx"
	"github.com/lcylpzls/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// TestNewValidation 覆盖构造校验。
func TestNewValidation(t *testing.T) {
	cases := []struct {
		name string
		cfg  Config
	}{
		{"空服务名", Config{}},
		{"负采样率", Config{ServiceName: "s", SampleRatio: -0.1}},
		{"超采样率", Config{ServiceName: "s", SampleRatio: 1.1}},
		{"未知导出器", Config{ServiceName: "s", Exporter: "unknown"}},
		{"OTLP 缺端点", Config{ServiceName: "s", Exporter: ExporterOTLPHTTP}},
	}
	for _, tc := range cases {
		if _, err := New(tc.cfg); !errx.Is(err, CodeInvalidConfig) {
			t.Fatalf("%s 应报配置错误，实际：%v", tc.name, err)
		}
	}
}

// TestExporterFailures 覆盖导出器构造失败分支。
func TestExporterFailures(t *testing.T) {
	origStdout := buildStdout
	buildStdout = func(io.Writer) (sdktrace.SpanExporter, error) {
		return nil, fmt.Errorf("stdout 构造失败")
	}
	if _, err := New(Config{ServiceName: "svc"}); !errx.Is(err, CodeExporterFailed) {
		t.Fatalf("stdout 构造失败应报错，实际：%v", err)
	}
	buildStdout = origStdout

	origOTLP := buildOTLP
	buildOTLP = func(context.Context, ...otlptracehttp.Option) (sdktrace.SpanExporter, error) {
		return nil, fmt.Errorf("otlp 构造失败")
	}
	if _, err := New(Config{ServiceName: "svc", Exporter: ExporterOTLPHTTP, OTLPEndpoint: "127.0.0.1:4318"}); !errx.Is(err, CodeExporterFailed) {
		t.Fatalf("otlp 构造失败应报错，实际：%v", err)
	}
	buildOTLP = origOTLP
}

// TestNewOTLP 覆盖 OTLP/HTTP 选项分支。
func TestNewOTLP(t *testing.T) {
	m, err := New(Config{
		ServiceName:  "svc",
		Exporter:     ExporterOTLPHTTP,
		OTLPEndpoint: "127.0.0.1:4318",
		OTLPInsecure: true,
		OTLPHeaders:  map[string]string{"Authorization": "Bearer x"},
		OTLPTimeout:  2 * time.Second,
	})
	testx.RequireNoError(t, err)

	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("OTLP 关闭失败：%v", err)
	}
}

// TestSampler 覆盖可插拔采样器。
func TestSampler(t *testing.T) {
	m, err := New(Config{
		ServiceName: "svc",
		Exporter:    ExporterMemory,
		Sampler:     sdktrace.NeverSample(),
	})
	testx.RequireNoError(t, err)

	_, span := m.Start(context.Background(), "never")
	span.End()
	_ = m.Shutdown(context.Background())
	if len(m.Spans()) != 0 {
		t.Fatal("NeverSample 不应导出任何 span")
	}
}

// TestSetGlobal 覆盖注册 OTel 全局组件。
func TestSetGlobal(t *testing.T) {
	prevProvider := otel.GetTracerProvider()
	prevProp := otel.GetTextMapPropagator()
	defer func() {
		otel.SetTracerProvider(prevProvider)
		otel.SetTextMapPropagator(prevProp)
	}()
	m, err := New(Config{
		ServiceName: "svc",
		Exporter:    ExporterMemory,
		SetGlobal:   true,
	})
	testx.RequireNoError(t, err)

	if otel.GetTracerProvider() == prevProvider {
		t.Fatal("全局 TracerProvider 未被替换")
	}
	ctx, span := m.Start(context.Background(), "global")
	carrier := propagation.HeaderCarrier(http.Header{})
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	span.End()
	if carrier.Get("traceparent") == "" {
		t.Fatal("全局传播器应能注入 traceparent")
	}
	_ = m.Shutdown(context.Background())
}

// TestConcurrentRequests 覆盖并发安全（race 检测）。
func TestConcurrentRequests(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	client := &http.Client{Transport: m.RoundTripper(http.DefaultTransport)}
	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			handler.ServeHTTP(httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, "http://example.com/in", nil))
			if _, err := client.Get(srv.URL + "/out"); err != nil {
				t.Errorf("出站请求失败：%v", err)
			}
		}()
	}
	wg.Wait()
	_ = m.Shutdown(context.Background())
	if got := len(m.Spans()); got != 2*n {
		t.Fatalf("应导出 %d 条 span，实际：%d", 2*n, got)
	}
}

// TestShutdownFailure 覆盖关闭失败分支。
func TestShutdownFailure(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	orig := shutdownProvider
	shutdownProvider = func(context.Context, *sdktrace.TracerProvider) error {
		return fmt.Errorf("关闭失败")
	}
	err = m.Shutdown(context.Background())
	shutdownProvider = orig
	if !errx.Is(err, CodeShutdownFailed) {
		t.Fatalf("关闭失败应报错，实际：%v", err)
	}
}

// TestNewStdout 覆盖 stdout 导出器。
func TestNewStdout(t *testing.T) {
	var buf bytes.Buffer
	m, err := New(Config{ServiceName: "svc", Writer: &buf})
	testx.RequireNoError(t, err)

	if m.Tracer() == nil || m.Propagator() == nil {
		t.Fatal("Tracer/Propagator 为空")
	}
	if m.Spans() != nil {
		t.Fatal("stdout 导出器不应提供内存快照")
	}
	ctx, span := m.Start(context.Background(), "op")
	span.End()
	_ = ctx
	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("关闭失败：%v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("stdout 导出器应输出内容")
	}
	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("重复关闭应幂等：%v", err)
	}
}

// TestNewMemory 覆盖内存导出器与 Span 快照。
func TestNewMemory(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	ctx, root := m.Start(context.Background(), "root",
		trace.WithAttributes(attribute.String("service.name", "svc")))
	_, child := m.Start(ctx, "child")
	child.End()
	root.End()
	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("关闭失败：%v", err)
	}
	spans := m.Spans()
	if len(spans) != 2 {
		t.Fatalf("应收集 2 条 Span，实际：%d", len(spans))
	}
	var rootSpan, childSpan SpanSnapshot
	for _, s := range spans {
		switch s.Name {
		case "root":
			rootSpan = s
		case "child":
			childSpan = s
		}
	}
	if rootSpan.Name == "" || childSpan.Name == "" {
		t.Fatalf("未找到父子 Span：%+v", spans)
	}
	if rootSpan.TraceID == "" || rootSpan.SpanID == "" || childSpan.ParentSpanID != rootSpan.SpanID {
		t.Fatalf("父子关系不符：root=%+v child=%+v", rootSpan, childSpan)
	}
	if rootSpan.Attributes["service.name"] != "svc" {
		t.Fatalf("属性不符：%+v", rootSpan.Attributes)
	}
	testx.RequireEqual(t, rootSpan.StatusCode, "Unset")

	m.mem.Reset()
	if len(m.Spans()) != 0 {
		t.Fatal("Reset 后应无 Span")
	}
}

// TestManagerMethods 覆盖配置副本、注入与提取。
func TestManagerMethods(t *testing.T) {
	cfg := Config{ServiceName: "svc", Version: "1.0.0", Environment: "test", Exporter: ExporterMemory}
	m, err := New(cfg)
	testx.RequireNoError(t, err)

	if got := m.Config(); got.ServiceName != "svc" || got.Version != "1.0.0" || got.Environment != "test" {
		t.Fatalf("配置副本不符：%+v", got)
	}
	parentCtx, span := m.Start(context.Background(), "parent")
	carrier := propagation.HeaderCarrier(http.Header{})
	m.Inject(parentCtx, carrier)
	extracted := m.Extract(context.Background(), carrier)
	sc := trace.SpanContextFromContext(extracted)
	if sc.TraceID() != trace.SpanContextFromContext(parentCtx).TraceID() {
		t.Fatal("注入/提取链路 ID 不一致")
	}
	span.End()
	_ = m.Shutdown(context.Background())
}

// TestMiddleware 覆盖状态码与错误标记。
func TestMiddleware(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	req := httptest.NewRequest(http.MethodGet, "http://example.com/hello", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("响应不符：%d %q", rec.Code, rec.Body.String())
	}

	errHandler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	errRec := httptest.NewRecorder()
	errHandler.ServeHTTP(errRec, httptest.NewRequest(http.MethodPost, "http://example.com/fail", nil))
	testx.RequireEqual(t, errRec.Code, http.StatusInternalServerError)

	_ = m.Shutdown(context.Background())

	spans := m.Spans()
	if len(spans) != 2 {
		t.Fatalf("应收集 2 条 Span，实际：%d", len(spans))
	}
	var okSpan, errSpan SpanSnapshot
	for _, s := range spans {
		switch s.Attributes["url.path"] {
		case "/hello":
			okSpan = s
		case "/fail":
			errSpan = s
		}
	}
	if okSpan.Name == "" || errSpan.Name == "" {
		t.Fatalf("未找到中间件 Span：%+v", spans)
	}
	if okSpan.Attributes["http.request.method"] != http.MethodGet ||
		okSpan.Attributes["url.path"] != "/hello" ||
		okSpan.Attributes["http.response.status_code"] != "200" {
		t.Fatalf("OK Span 属性不符：%+v", okSpan.Attributes)
	}
	if errSpan.StatusCode != "Error" || errSpan.StatusMessage != "HTTP 500" {
		t.Fatalf("错误 Span 状态不符：%+v", errSpan)
	}
}

// TestMiddlewareRouteNaming 覆盖 ServeMux 路由模板命名。
func TestMiddlewareRouteNaming(t *testing.T) {
	m, err := New(Config{
		ServiceName: "svc",
		Exporter:    ExporterMemory,
		RouteNamer: func(r *http.Request) string {
			return "/users/{id}"
		},
	})
	testx.RequireNoError(t, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	handler.ServeHTTP(httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "http://example.com/users/42", nil))
	_ = m.Shutdown(context.Background())

	found := false
	for _, s := range m.Spans() {
		if s.Name == "GET /users/{id}" {
			found = true
		}
	}
	testx.RequireTrue(t, found)

	// RouteNamer 返回空串时保持默认命名。
	m2, err := New(Config{
		ServiceName: "svc",
		Exporter:    ExporterMemory,
		RouteNamer:  func(*http.Request) string { return "" },
	})
	testx.RequireNoError(t, err)

	h2 := m2.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	h2.ServeHTTP(httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "http://example.com/plain", nil))
	_ = m2.Shutdown(context.Background())
	for _, s := range m2.Spans() {
		if s.Name == "/plain" {
			return
		}
	}
	t.Fatalf("空路由名应保持默认命名：%+v", m2.Spans())
}

// TestMiddlewareSlow 覆盖慢请求标记。
func TestMiddlewareSlow(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory, SlowThreshold: time.Nanosecond})
	testx.RequireNoError(t, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "http://example.com/slow", nil))
	_ = m.Shutdown(context.Background())

	for _, s := range m.Spans() {
		if s.Name == "/slow" {
			if s.Attributes["request.duration_ms"] == "" {
				t.Fatalf("慢请求应记录耗时属性：%+v", s.Attributes)
			}
			for _, ev := range s.Events {
				if ev.Name == "slow" {
					return
				}
			}
			t.Fatalf("慢请求应记录 slow 事件：%+v", s.Events)
		}
	}
	t.Fatal("未找到慢请求 span")
}

// TestMiddlewarePropagation 覆盖入站链路上下文透传。
func TestMiddlewarePropagation(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	parentCtx, parent := m.Start(context.Background(), "parent")
	sc := trace.SpanContextFromContext(parentCtx)
	header := http.Header{}
	header.Set("traceparent", fmt.Sprintf("00-%s-%s-01", sc.TraceID(), sc.SpanID()))

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	req := httptest.NewRequest(http.MethodGet, "http://example.com/child", nil)
	req.Header = header
	handler.ServeHTTP(httptest.NewRecorder(), req)
	parent.End()
	_ = m.Shutdown(context.Background())

	for _, s := range m.Spans() {
		if s.Name == "/child" {
			if s.ParentSpanID != sc.SpanID().String() {
				t.Fatalf("入站 Span 父 ID 不符：%s != %s", s.ParentSpanID, sc.SpanID())
			}
			return
		}
	}
	t.Fatal("未找到入站 Span")
}

// TestLogFields 覆盖日志字段生成。
func TestLogFields(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	ctx, span := m.Start(context.Background(), "op")
	sc := trace.SpanContextFromContext(ctx)
	span.End()
	_ = m.Shutdown(context.Background())

	var buf bytes.Buffer
	logger, err := logx.NewBuilder().EnableWriter(&buf, logx.InfoLevel).Build()
	testx.RequireNoError(t, err)

	logger.Info("带链路", LogFields(ctx))
	if !strings.Contains(buf.String(), sc.TraceID().String()) {
		t.Fatalf("日志应包含 trace_id：%s", buf.String())
	}

	buf.Reset()
	logger.Info("无链路", LogFields(context.Background()))
	if strings.Contains(buf.String(), "trace_id") {
		t.Fatalf("无链路不应输出 trace_id：%s", buf.String())
	}
}

// TestMemoryShutdown 覆盖导出器关闭后忽略导出。
func TestMemoryShutdown(t *testing.T) {
	m, err := New(Config{ServiceName: "svc", Exporter: ExporterMemory})
	testx.RequireNoError(t, err)

	_ = m.mem.Shutdown(context.Background())
	_, span := m.Start(context.Background(), "after-shutdown")
	span.End()
	_ = m.Shutdown(context.Background())
	if len(m.Spans()) != 0 {
		t.Fatal("关闭后不应再收集 Span")
	}
}
