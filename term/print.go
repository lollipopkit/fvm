package term

import (
	"os"
)

const (
	red     = "\033[91m"
	green   = "\033[92m"
	yellow  = "\033[93m"
	blue    = "\033[94m"
	cyan    = "\033[96m"
	white   = "\033[97m"
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
