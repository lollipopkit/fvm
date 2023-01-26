package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "use",
		Aliases:   []string{"u"},
		Usage:     "Use a specific version of Flutter",
		Action:    handleUse,
		ArgsUsage: "[version]",
		UsageText: consts.APP_NAME + " use [version]",
	})
}

func handleUse(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() != 1 {
		term.Warn("Usage: " + ctx.Command.UsageText)
	}

	version := args.Get(0)
	err := utils.Use(version)
	if err != nil {
		if err == utils.ErrVersionNotInstalled {
			term.Warn(err.Error())
			confirm := term.Confirm("Do you want to install it?", true)
			if confirm {
				releases, err := utils.GetReleases()
				if err != nil {
					return err
				}

				release, err := utils.GetReleaseByVersion(releases ,version)
				if err != nil {
					return err
				}

				return utils.Install(release)
			}
		} else {
			return err
		}
	}

	return nil
}
