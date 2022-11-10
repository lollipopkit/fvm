package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "config",
		Aliases:   []string{"c"},
		Usage:     "Config something",
		Subcommands: []*cli.Command{
			{
				Name:      "alias",
				Aliases:   []string{"a"},
				Usage:     "Config alias `dart -> fvm dart` and `flutter -> fvm flutter`",
				Action:    handleConfigAlias,
				ArgsUsage: "[alias]",
				UsageText: consts.APP_NAME + " config alias [alias]",
			},
		},
	})
}

func handleConfigAlias(ctx *cli.Context) error {
	return utils.SetAlias()
}
