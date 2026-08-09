package tracex

import (
	"testing"

	"github.com/lcylpzls/errx"
)

// TestCodesRegistered 覆盖错误码注册与匹配。
func TestCodesRegistered(t *testing.T) {
	for _, code := range []errx.Code{CodeInvalidConfig, CodeExporterFailed, CodeShutdownFailed} {
		e := errx.NewCode(code, "测试错误")
		if !errx.Is(e, code) {
			t.Fatalf("错误码 %s 无法匹配", code)
		}
		if e.Kind() == errx.KindUnknown {
			t.Fatalf("错误码 %s 未注册分类", code)
		}
	}
}
