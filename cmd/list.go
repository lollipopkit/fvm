package cmd

import (
	"io/ioutil"

	"github.com/lollipopkit/gofvm/consts"
	"github.com/lollipopkit/gofvm/term"
	"github.com/lollipopkit/gofvm/utils"
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
	p := utils.Path()
	dirs, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}

	term.Cyan("Installed versions:")
	for _, dir := range dirs {
		if dir.IsDir() {
			println(dir.Name())
		}
	}

	return nil
}
