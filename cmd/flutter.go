package cmd

import (
	"os"
	"path/filepath"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:            "flutter",
		Aliases:         []string{"f"},
		Usage:           "Proxy flutter commands",
		Action:          handleFlutter,
		SkipFlagParsing: true,
	})
}

func handleFlutter(ctx *cli.Context) error {
	args := ctx.Args().Slice()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	wdFvm := filepath.Join(pwd, consts.FvmDirName)
	if utils.Exists(wdFvm) {
		err = utils.Execute(filepath.Join(wdFvm, "bin/flutter"), args...)
	} else {
		err = utils.Execute("flutter", args...)
	}
	return err
}
