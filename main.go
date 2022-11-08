package main

import (
	"github.com/LollipopKit/gofvm/cmd"
	"github.com/LollipopKit/gofvm/term"
	"github.com/LollipopKit/gofvm/utils"
)

func main() {
	err := utils.Precheck()
	if err != nil {
		term.Yellow(err.Error())
		return
	}

	cmd.Run()
}
