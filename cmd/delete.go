package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "delete",
		Aliases:   []string{"d"},
		Usage:     "Delete a specific version of Flutter",
		Action:    handleDelete,
		ArgsUsage: "[version]",
		UsageText: consts.APP_NAME + " delete [version]",
	})
}

func handleDelete(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		term.Warn("Usage: " + ctx.Command.UsageText)
	}

	version := ctx.Args().Get(0)
	if !utils.IsVersionInstalled(version) {
		term.Warn("Version [" + version + "] is not installed")
		return nil
	}

	err := utils.Delete(version)
	if err != nil {
		term.Warn("Failed to delete version [" + version + "]")
		return err
	}

	term.Success("Deleted version [" + version + "]")
	return nil
}
