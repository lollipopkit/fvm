package term

import (
	"os"
)

const (
	red     = "\033[91m"
	green   = "\033[32m"
	yellow  = "\033[93m"
	cyan    = "\033[96m"
	noColor = "\033[0m"
)

func print(s string) {
	os.Stdout.WriteString(s + noColor + "\n")
}

func Error(s string, noPanic ...bool) {
	if len(noPanic) > 0 && noPanic[0] {
		print(red + s)
	} else {
		panic(s)
	}
}

func Success(s string) {
	print(green + s)
}

func Warn(s string) {
	print(yellow + s)
}

func Info(s string) {
	print(cyan + s)
}
