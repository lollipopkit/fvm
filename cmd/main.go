package cmd

import (
	"os"

	"github.com/LollipopKit/gofvm/consts"
	"github.com/LollipopKit/gofvm/term"
	"github.com/urfave/cli/v2"
)

var (
	cmds  = []*cli.Command{}
	flags = []cli.Flag{}
)

func Run() {
	app := &cli.App{
		Name:     consts.APP_NAME,
		Usage:    "Flutter Version Manager written in Go",
		Commands: cmds,
		Flags:    flags,
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() == 0 {
				return cli.ShowAppHelp(ctx)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		term.Red(err.Error(), true)
	}
}
