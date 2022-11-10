package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "global",
		Aliases:   []string{"g"},
		Usage:     "Manage global version of Flutter",
		UsageText: consts.APP_NAME + " global [version]",
		ArgsUsage: "[version]",
		Action:    handleGlobal,
	})
}

func handleGlobal(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() != 1 {
		term.Warn("Usage: " + ctx.Command.UsageText)
	}

	err := utils.Global(args.Get(0))
	if err != nil {
		if err == utils.ErrVersionNotInstalled {
			term.Warn(err.Error())
		} else {
			return err
		}
	}

	return nil
}
