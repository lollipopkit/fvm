package cmd

import (
	"os"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/gommon/log"
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
		Version:  consts.APP_VERSION,
		Commands: cmds,
		Flags:    flags,
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() == 0 {
				return cli.ShowAppHelp(ctx)
			}
			return nil
		},
		Suggest:              true,
		EnableBashCompletion: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err.Error(), true)
	}
}
