package cmd

import (
	"os"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
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
		Suggest: true,
	}

	if err := app.Run(os.Args); err != nil {
		term.Error(err.Error(), true)
	}
}
