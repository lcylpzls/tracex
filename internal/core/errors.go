// 错误码定义：统一以 TRACEX_ 为前缀（errx 家族规范）。
package core

import "github.com/lcylpzls/errx"

const (
	CodeInvalidConfig  errx.Code = "TRACEX_INVALID_CONFIG"
	CodeExporterFailed errx.Code = "TRACEX_EXPORTER_FAILED"
	CodeShutdownFailed errx.Code = "TRACEX_SHUTDOWN_FAILED"
)

func init() {
	errx.RegisterCode(CodeInvalidConfig, "配置非法")
	errx.RegisterCodeKind(CodeInvalidConfig, errx.KindInvalid)
	errx.RegisterCode(CodeExporterFailed, "创建追踪导出器失败")
	errx.RegisterCodeKind(CodeExporterFailed, errx.KindUnavailable)
	errx.RegisterCode(CodeShutdownFailed, "关闭追踪器失败")
	errx.RegisterCodeKind(CodeShutdownFailed, errx.KindUnavailable)
}
