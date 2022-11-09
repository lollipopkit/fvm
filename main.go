package main

import (
	"github.com/lollipopkit/gofvm/cmd"
	"github.com/lollipopkit/gofvm/term"
	"github.com/lollipopkit/gofvm/utils"
)

func main() {
	err := utils.Precheck()
	if err != nil {
		term.Yellow(err.Error())
		return
	}

	cmd.Run()
}
