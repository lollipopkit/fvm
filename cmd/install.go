package cmd

import (
	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "install",
		Aliases:   []string{"i"},
		Usage:     "Install a specific version of Flutter",
		Action:    handleInstall,
		ArgsUsage: "[version]",
		UsageText: consts.APP_NAME + " install [version]",
	})
}

func handleInstall(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		term.Warn("Usage: " + ctx.Command.UsageText)
	}

	releases, err := utils.GetReleases()
	if err != nil {
		return err
	}
	vs := make([]string, 0)
	for idx := range releases {
		vs = append(vs, releases[idx].Version)
	}

	version := ctx.Args().Get(0)
	if !utils.Contains(vs, version) {
		term.Warn("Version [" + version + "] is not available")
		return nil
	}

	r, err := utils.GetReleaseByVersion(releases, version)
	if err != nil {
		return err
	}
	err = utils.Install(r)
	if err != nil {
		return err
	}

	return nil
}
