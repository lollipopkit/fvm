package cmd

import (
	"os"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/lollipopkit/gommon/term"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "list",
		Aliases:   []string{"l"},
		Usage:     "List all installed versions of Flutter",
		UsageText: consts.APP_NAME + " list",
		Action:    handleList,
	})
}

func handleList(ctx *cli.Context) error {
	dirs, err := os.ReadDir(utils.FvmHome)
	if err != nil {
		return err
	}

	gVersion, err := utils.GetGlobalVersion()
	if err != nil {
		term.Warn("You have not set a global version yet.")
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			dName := dir.Name()
			if dName == gVersion {
				term.Yellow(dName + " [GLOBAL]")
			} else {
				println(dName)
			}
		}
	}

	return nil
}
