package core

import (
	"testing"

	"github.com/lcylpzls/errx"
	"github.com/lcylpzls/testx"
)

// TestCodesRegistered 覆盖错误码注册与匹配。
func TestCodesRegistered(t *testing.T) {
	for _, code := range []errx.Code{CodeInvalidConfig, CodeExporterFailed, CodeShutdownFailed} {
		e := errx.NewCode(code, "测试错误")
		testx.RequireErrCode(t, e, code)
		testx.RequireNotEqual(t, e.Kind(), errx.KindUnknown)
	}
}
