package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/lollipopkit/gommon/term"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Config something",
		Subcommands: []*cli.Command{
			{
				Name:      "alias",
				Aliases:   []string{"a"},
				Usage:     "Config alias `dart -> fvm dart` and `flutter -> fvm flutter`",
				Action:    handleConfigAlias,
				ArgsUsage: "[alias]",
				UsageText: consts.APP_NAME + " config alias [alias]",
			},
			{
				Name:      "use-mirror",
				Aliases:   []string{"um"},
				Usage:     "config use mirror or not",
				Action:    handleConfigMirror,
				UsageText: consts.APP_NAME + " config use-mirror [ true | false ]",
			},
		},
	})
}

func handleConfigAlias(ctx *cli.Context) error {
	return utils.SetAlias()
}

func handleConfigMirror(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() != 1 {
		term.Warn("Usage: " + ctx.Command.UsageText)
		return nil
	}
	use := args.Get(0)
	if use == "true" {
		a := true
		utils.Config.UseMirror = &a
		return utils.SaveConfig()
	} else if use == "false" {
		a := false
		utils.Config.UseMirror = &a
		return utils.SaveConfig()
	}
	term.Warn("Usage: " + ctx.Command.UsageText)
	return nil
}
