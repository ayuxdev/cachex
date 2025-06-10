package logger

import "github.com/fatih/color"

var (
	errorColor     = color.New(color.FgRed).Add(color.Bold)
	infoColor      = color.New(color.FgCyan).Add(color.Bold)
	successColor   = color.New(color.FgGreen).Add(color.Bold)
	warnColor      = color.New(color.FgYellow).Add(color.Bold)
	debugColor     = color.New(color.FgBlue).Add(color.Bold)
	timeStampColor = color.New(color.FgMagenta).Add(color.Bold)
)
