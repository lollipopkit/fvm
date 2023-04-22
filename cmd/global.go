package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/lollipopkit/gommon/term"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "global",
		Aliases:   []string{"g"},
		Usage:     "Manage global version of Flutter",
		UsageText: consts.APP_NAME + " global [version]",
		ArgsUsage: "[version]",
		Action:    handleGlobal,
	})
}

func handleGlobal(ctx *cli.Context) error {
	args := ctx.Args()
	argsLen := args.Len()
	if argsLen == 0 {
		ver, err := utils.GetGlobalVersion()
		if err != nil {
			term.Err(err.Error())
			term.Info("Usage: " + ctx.Command.UsageText)
			return nil
		}
		term.Info("Global version: " + ver)
		return nil
	}
	if argsLen > 1 {
		term.Warn("Usage: " + ctx.Command.UsageText)
		return nil
	}

	ver := args.Get(0)
	err := utils.Global(ver)
	if err != nil {
		if err == utils.ErrVersionNotInstalled {
			term.Warn(err.Error())
			confirm := term.Confirm("Install " + ver + " now?", true)
			if confirm {
				return install(ver)
			}
		} else {
			return err
		}
	}

	return nil
}
