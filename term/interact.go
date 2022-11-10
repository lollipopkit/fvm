package term

import (
	"fmt"
	"strconv"
	"strings"
)

func ReadLine() string {
	var input string
	fmt.Scanln(&input)
	return input
}

func Confirm(question string, default_ bool) bool {
	suffix := func() string {
		if default_ {
			return " [Y/n]"
		}
		return " [y/N]"
	}()
	Info(fmt.Sprintf("%s%s: ", question, suffix))
	input := ReadLine()
	if input == "" {
		return default_
	}
	return strings.ToLower(input) == "y"
}

func Option(question string, options []string, default_ int) int {
	println()
	for i := range options {
		print(fmt.Sprintf("%d. %s ", i+1, options[i]))
	}
	suffix := fmt.Sprintf(" [default: %d]", default_+1)
	Info(fmt.Sprintf("%s %s:", question, suffix))
	input := ReadLine()
	if input == "" {
		return default_
	}
	inputIdx, err := strconv.Atoi(input)
	if err != nil {
		return default_
	}
	return inputIdx - 1
}
