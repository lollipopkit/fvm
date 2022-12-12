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
		Name:            "dart",
		Aliases:         []string{"d"},
		Usage:           "Proxy dart commands",
		Action:          handleDart,
		SkipFlagParsing: true,
	})
}

func handleDart(ctx *cli.Context) error {
	args := ctx.Args().Slice()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	wdFvm := filepath.Join(pwd, consts.FvmDirName)
	if utils.Exists(wdFvm) {
		err = utils.Execute(filepath.Join(wdFvm, "bin/dart"), args...)
	} else {
		err = utils.Execute("dart", args...)
	}
	return err
}
