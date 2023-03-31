package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/utils"
	"github.com/lollipopkit/gommon/term"
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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    consts.FlagForce,
				Aliases: []string{"f"},
				Usage:   "Force use (download if not installed))",
			},
		},
	})
}

func handleUse(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() != 1 {
		term.Warn("Usage: " + ctx.Command.UsageText)
		return nil
	}

	version := args.Get(0)
	err := utils.Use(version)
	if err != nil {
		if err == utils.ErrVersionNotInstalled {
			term.Warn(err.Error() + "\n")
			if ctx.Bool(consts.FlagForce) || term.Confirm("Do you want to install it?", true) {
				return install(version)
			}
		} else {
			return err
		}
	}
	return nil
}

func install(version string) error {
	releases, err := utils.GetReleases()
	if err != nil {
		return err
	}

	release, err := utils.GetReleaseByVersion(releases, version)
	if err != nil {
		return err
	}

	err = utils.Install(release, false)
	if err != nil {
		return err
	}
	return utils.Use(version)
}
