package term

import (
	"os"
)

const (
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	noColor = "\033[0m"
)

func print(s string) {
	os.Stdout.WriteString(s + noColor + "\n")
}

func Red(s string, noPanic ...bool) {
	if len(noPanic) > 0 && noPanic[0] {
		print(red + s)
	} else {
		panic(s)
	}
}

func Green(s string) {
	print(green + s)
}

func Yellow(s string) {
	print(yellow + s)
}

func Blue(s string) {
	print(blue + s)
}

func Cyan(s string) {
	print(cyan + s)
}

func White(s string) {
	print(white + s)
}
