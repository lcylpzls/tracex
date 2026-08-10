package core

import (
	"context"
	testx "github.com/lcylpzls/testx"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkMiddleware 基准：HTTP 中间件全链路开销。
func BenchmarkMiddleware(b *testing.B) {
	m, err := New(Config{ServiceName: "bench", Exporter: ExporterMemory})
	testx.RequireNoError(b, err)

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "http://example.com/bench", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}
	b.StopTimer()
	_ = m.Shutdown(context.Background())
}

// BenchmarkRoundTripper 基准：出站注入开销。
func BenchmarkRoundTripper(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	m, err := New(Config{ServiceName: "bench", Exporter: ExporterMemory})
	testx.RequireNoError(b, err)

	client := &http.Client{Transport: m.RoundTripper(http.DefaultTransport)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := client.Get(srv.URL); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	_ = m.Shutdown(context.Background())
}
