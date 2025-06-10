package cachex

import (
	"github.com/ayuxdev/cachex/internal/pkg/logger"
	"github.com/ayuxdev/cachex/internal/scanner"
)

// mapLogMode maps the string log mode to the corresponding scanner.LogMode
func mapStrLogMode(mode string) scanner.LogMode {
	switch mode {
	case "pretty":
		return scanner.PrettyLog
	case "json":
		return scanner.JsonLog
	case "":
		return scanner.PrettyLog
	default:
		logger.Errorf("invalid log mode: %s, defaulting to pretty", mode)
		return scanner.PrettyLog
	}
}

// mapLogTarget maps the string log target to the corresponding scanner.LogTarget
func mapStrLogTarget(target string) scanner.LogTarget {
	switch target {
	case "stdout":
		return scanner.StdoutLog
	case "file":
		return scanner.FileLog
	case "both":
		return scanner.BothLog
	case "":
		return scanner.StdoutLog
	default:
		logger.Errorf("invalid log target: %s, defaulting to stdout", target)
		return scanner.StdoutLog
	}
}

// mapScanMode maps the string scan mode to the corresponding scanner.ScanMode
func mapStrScanMode(scanMode string) scanner.ScanMode {
	switch scanMode {
	case "single":
		return scanner.SingleHeaderScanMode
	case "multi":
		return scanner.MultiHeaderScanMode
	case "":
		return scanner.SingleHeaderScanMode
	default:
		logger.Errorf("invalid scan mode: %s, defaulting to single header scan", scanMode)
		return scanner.SingleHeaderScanMode
	}
}
