// demo 演示 tracex：stdout 导出器 + HTTP 中间件 + 出站注入 + 日志钩子。
package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lcylpzls/logx"
	"github.com/lcylpzls/tracex"
)

func main() {
	m, err := tracex.New(tracex.Config{
		ServiceName: "tracex-demo",
		Environment: "dev",
	})
	if err != nil {
		panic(err)
	}
	defer m.Shutdown(context.Background())

	logger, err := logx.NewBuilder().
		EnableWriter(io.Discard, logx.InfoLevel).
		Build()
	if err != nil {
		panic(err)
	}
	if hl, ok := logger.(logx.HookedLogger); ok {
		hl.AddHook(tracex.NewLogHook())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		logger.WithContext(r.Context()).Info("收到请求", logx.Fields())
		client := &http.Client{Transport: m.RoundTripper(http.DefaultTransport)}
		_, _ = client.Get("https://example.com")
		_, _ = w.Write([]byte("hello"))
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		_ = http.ListenAndServe("127.0.0.1:8080", m.Middleware(mux))
	}()
	<-stop
}
