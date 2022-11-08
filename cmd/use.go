package cmd

import (
	"github.com/LollipopKit/gofvm/consts"
	"github.com/LollipopKit/gofvm/term"
	"github.com/LollipopKit/gofvm/utils"
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
	if ctx.Args().Len() != 1 {
		term.Yellow("Usage: " + ctx.Command.UsageText)
	}

	err := utils.Use(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	return nil
}
