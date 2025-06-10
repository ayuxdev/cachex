package logger

import (
	"fmt"
	"os"
	"time"
)

var (
	DisableWarn     = false
	DisableInfo     = false
	DisableDebug    = false
	DisableSuccess  = false
	EnableTimeStamp = false
)

func formatMessage(msg string) string {
	if EnableTimeStamp {
		timeStamp := time.Now().Format("2006-01-02 15:04:05")
		return timeStampColor.Sprintf("[%s]", timeStamp) + fmt.Sprintf(" %s", msg)
	}
	return msg
}

func Error(message string) {
	errorColor.Fprint(os.Stderr, "[error] ")
	fmt.Fprintln(os.Stderr, formatMessage(message))
}

func Errorf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}

func Info(message string) {
	if DisableInfo {
		return
	}
	infoColor.Fprint(os.Stderr, "[info] ")
	fmt.Fprintln(os.Stderr, formatMessage(message))
}

func Infof(format string, a ...any) {
	Info(fmt.Sprintf(format, a...))
}

func Success(message string) {
	if DisableSuccess {
		return
	}
	successColor.Fprint(os.Stderr, "[ok] ")
	fmt.Fprintln(os.Stderr, formatMessage(message))
}

func Successf(format string, a ...any) {
	Success(fmt.Sprintf(format, a...))
}

func Warn(message string) {
	if DisableWarn {
		return
	}
	warnColor.Fprint(os.Stderr, "[warn] ")
	fmt.Fprintln(os.Stderr, formatMessage(message))
}

func Warnf(format string, a ...any) {
	Warn(fmt.Sprintf(format, a...))
}

func Debug(message string) {
	if DisableDebug {
		return
	}
	debugColor.Fprint(os.Stderr, "[debug] ")
	fmt.Fprintln(os.Stderr, formatMessage(message))
}

func Debugf(format string, a ...any) {
	Debug(fmt.Sprintf(format, a...))
}
