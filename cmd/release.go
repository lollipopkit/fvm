package cmd

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/model"
	"github.com/lollipopkit/fvm/utils"
	"github.com/lollipopkit/gommon/term"
	"github.com/urfave/cli/v2"
)

var (
	majorVersionReg = regexp.MustCompile(`^v?(\d+)\.\S+$`)
	// stable: X.X.X
	stableVersionReg = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "release",
		Aliases:   []string{"r"},
		Usage:     "List all releases of Flutter",
		Action:    handleRelease,
		UsageText: consts.APP_NAME + " release",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "stable",
				Aliases: []string{"s"},
				Usage:   "List all stable releases of Flutter",
			},
			&cli.BoolFlag{
				Name:    "preview",
				Aliases: []string{"p", "pre"},
				Usage:   "List all preview releases of Flutter",
			},
		},
	})
}

func handleRelease(ctx *cli.Context) error {
	onlyStable := ctx.Bool("stable")
	onlyPreview := ctx.Bool("preview")
	if onlyStable && onlyPreview {
		return fmt.Errorf("can not use --stable and --preview at the same time")
	}
	releases, err := utils.GetReleases()
	if err != nil {
		return err
	}

	majorVersionsMap := make(map[string][]model.Release, 0)
	for idx := range releases {
		if onlyStable && !stableVersionReg.MatchString(releases[idx].Version) {
			continue
		}
		if onlyPreview && stableVersionReg.MatchString(releases[idx].Version) {
			continue
		}
		m := majorVersionReg.FindStringSubmatch(releases[idx].Version)
		if len(m) < 2 {
			continue
		}
		majorVersion := m[1]
		if _, ok := majorVersionsMap[majorVersion]; !ok {
			majorVersionsMap[majorVersion] = []model.Release{releases[idx]}
		} else {
			majorVersionsMap[majorVersion] = append(majorVersionsMap[majorVersion], releases[idx])
		}
	}

	majorVersions := make([]string, 0)
	for k := range majorVersionsMap {
		majorVersions = append(majorVersions, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(majorVersions)))

	arch := utils.GetArch()
	for _, majorVersion := range majorVersions {
		term.Green(fmt.Sprintf("[%s.x]:\n", majorVersion))
		count := 0
		printText := ""
		for _, release := range majorVersionsMap[majorVersion] {
			if release.DartSdkArch != arch {
				continue
			}
			count++
			if count > 5 {
				break
			}
			printText += fmt.Sprintf("%s\t%s\t%s\n", release.Version, release.DartSdkVersion, release.ReleaseDate.Format("2006-01-02"))
		}

		print(printText)
		if count > 5 {
			term.Yellow(fmt.Sprintf("...and %d more", len(majorVersionsMap[majorVersion])-count))
		}
		println()
		break
	}

	return nil
}
