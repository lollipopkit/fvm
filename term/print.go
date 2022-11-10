package term

const (
	red     = "\033[91m"
	green   = "\033[32m"
	yellow  = "\033[93m"
	cyan    = "\033[96m"
	noColor = "\033[0m"
)

func Error(s string, noPanic ...bool) {
	if len(noPanic) > 0 && noPanic[0] {
		print(red + s + noColor + "\n")
	} else {
		panic(s)
	}
}

func Success(s string) {
	print(green + s + noColor + "\n")
}

func SuccessNln(s string) {
	print(green + s + noColor)
}

func Warn(s string) {
	print(yellow + s + noColor + "\n")
}

func WarnNln(s string) {
	print(yellow + s + noColor)
}

func Info(s string) {
	print(cyan + s + noColor + "\n")
}

func InfoNln(s string) {
	print(cyan + s + noColor)
}
