package cmd

import (
	"os"
	"path"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "dart",
		Aliases:   []string{"d"},
		Usage:     "Proxy dart commands",
		UsageText: consts.APP_NAME + " dart [command]",
		Action:    handleDart,
	})
}

func handleDart(ctx *cli.Context) error {
	args := ctx.Args().Slice()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	wdFvm := path.Join(pwd, consts.FVM_DIR_NAME)
	if utils.Exists(wdFvm) {
		err = utils.Execute(path.Join(wdFvm, "bin/dart"), args...)
	} else {
		err = utils.Execute("dart", args...)
	}
	return err
}
