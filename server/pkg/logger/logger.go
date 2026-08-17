package logger

import "fmt"

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

func Success(format string, a ...any) {
	fmt.Printf("%s[SUCCESS]%s %s\n", green, reset, fmt.Sprintf(format, a...))
}

func Info(format string, a ...any) {
	fmt.Printf("%s[INFO]%s %s\n", blue, reset, fmt.Sprintf(format, a...))
}

func Warn(format string, a ...any) {
	fmt.Printf("%s[WARN]%s %s\n", yellow, reset, fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) {
	fmt.Printf("%s[ERROR]%s %s\n", red, reset, fmt.Sprintf(format, a...))
}
