package webx

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	testx "github.com/lcylpzls/testx"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lcylpzls/logx"
	"github.com/lcylpzls/tracex"
	wx "github.com/lcylpzls/webx"
)

// testLogger 构造丢弃日志器。
func testLogger() logx.Logger {
	l, err := logx.NewBuilder().EnableWriter(io.Discard, logx.InfoLevel).Build()
	if err != nil {
		panic(err)
	}
	return l
}

// writeTestCert 生成自签名证书。
func writeTestCert(t *testing.T) (string, string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	testx.RequireNoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "tracex-webx"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	testx.RequireNoError(t, err)

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, err := x509.MarshalECPrivateKey(key)
	testx.RequireNoError(t, err)

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")
	if err := os.WriteFile(certFile, certPEM, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, keyPEM, 0o600); err != nil {
		t.Fatal(err)
	}
	return certFile, keyFile
}

// TestMiddleware 覆盖 webx 中间件集成。
func TestMiddleware(t *testing.T) {
	certFile, keyFile := writeTestCert(t)
	m, err := tracex.New(tracex.Config{
		ServiceName:   "svc",
		Exporter:      tracex.ExporterMemory,
		SlowThreshold: time.Millisecond,
	})
	testx.RequireNoError(t, err)

	s := wx.NewServer(wx.Config{
		TLSCertFile:     certFile,
		TLSKeyFile:      keyFile,
		ShutdownTimeout: 5 * time.Second,
	}, testLogger())
	s.UseHttp2Listen("127.0.0.1:0")
	s.UseGlobalMiddleware(Middleware(m))
	s.RegisterRoute(wx.Route{
		Method: http.MethodGet,
		Path:   "/hello",
		Handler: func(c *wx.Context) {
			time.Sleep(5 * time.Millisecond)
			_ = c.String(http.StatusOK, "ok")
		},
	})
	s.RegisterRoute(wx.Route{
		Method: http.MethodGet,
		Path:   "/error",
		Handler: func(c *wx.Context) {
			_ = c.String(http.StatusInternalServerError, "boom")
		},
	})
	errCh := make(chan error, 1)
	go func() { errCh <- s.Start() }()
	var addr string
	for i := 0; i < 500; i++ {
		if addr = s.ListenerAddr(); addr != "" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	testx.RequireNotEqual(t, addr, "")

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Get("https://" + addr + "/hello")
	testx.RequireNoError(t, err)

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	testx.RequireEqual(t, resp.StatusCode, http.StatusOK)

	resp2, err := client.Get("https://" + addr + "/error")
	if err != nil || resp2.StatusCode != http.StatusInternalServerError {
		t.Fatalf("错误路由不符：%v %d", err, resp2.StatusCode)
	}
	_, _ = io.Copy(io.Discard, resp2.Body)
	_ = resp2.Body.Close()
	resp3, err := client.Get("https://" + addr + "/missing")
	if err != nil || resp3.StatusCode != http.StatusNotFound {
		t.Fatalf("缺失路由不符：%v %d", err, resp3.StatusCode)
	}
	_, _ = io.Copy(io.Discard, resp3.Body)
	_ = resp3.Body.Close()
	_ = s.Stop(context.Background())
	_ = m.Shutdown(context.Background())

	var okSpan, errSpan, missingSpan tracex.SpanSnapshot
	for _, sp := range m.Spans() {
		switch sp.Name {
		case "GET /hello":
			okSpan = sp
		case "GET /error":
			errSpan = sp
		case "/missing":
			missingSpan = sp
		}
	}
	if okSpan.Name == "" || okSpan.Attributes["http.response.status_code"] != "200" {
		t.Fatalf("未找到 webx 链路 span：%+v", m.Spans())
	}
	slowFound := false
	for _, ev := range okSpan.Events {
		if ev.Name == "slow" {
			slowFound = true
		}
	}
	testx.RequireTrue(t, slowFound)

	if errSpan.StatusCode != "Error" || errSpan.StatusMessage != "HTTP 500" {
		t.Fatalf("5xx span 状态不符：%+v", errSpan)
	}
	if missingSpan.Name != "/missing" || missingSpan.Attributes["http.response.status_code"] != "404" {
		t.Fatalf("404 span 应被全局中间件覆盖：%+v", missingSpan)
	}
}
