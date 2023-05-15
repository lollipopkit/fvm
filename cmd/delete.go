package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/lollipopkit/gommon/log"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "delete",
		Aliases:   []string{"D"},
		Usage:     "Delete a specific version of Flutter",
		Action:    handleDelete,
		ArgsUsage: "[version]",
		UsageText: consts.APP_NAME + " delete [version]",
	})
}

func handleDelete(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		log.Warn("Usage: " + ctx.Command.UsageText)
		return nil
	}

	version := ctx.Args().Get(0)
	if !utils.IsVersionInstalled(version) {
		log.Warn("Version [" + version + "] is not installed")
		return nil
	}

	err := utils.Delete(version)
	if err != nil {
		log.Warn("Failed to delete version [" + version + "]")
		return err
	}

	log.Suc("Deleted version [" + version + "]")
	return nil
}
