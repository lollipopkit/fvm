package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "use",
		Aliases:   []string{"u"},
		Usage:     "Use a specific version of Flutter",
		Action:    handleUse,
		ArgsUsage: "[version]",
		UsageText: consts.APP_NAME + " use [version]",
	})
}

func handleUse(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() != 1 {
		term.Yellow("Usage: " + ctx.Command.UsageText)
	}

	err := utils.Use(args.Get(0))
	if err != nil {
		if err == utils.ErrVersionNotInstalled {
			term.Yellow(err.Error())
		} else {
			return err
		}
	}

	return nil
}
