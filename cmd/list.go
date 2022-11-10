package cmd

import (
	"io/ioutil"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/lollipopkit/fvm/utils"
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
	dirs, err := ioutil.ReadDir(utils.FvmHome)
	if err != nil {
		return err
	}

	term.Info("Installed versions:")
	for _, dir := range dirs {
		if dir.IsDir() {
			println(dir.Name())
		}
	}

	return nil
}
