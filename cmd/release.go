package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/model"
	"github.com/lollipopkit/fvm/term"
	"github.com/lollipopkit/fvm/utils"
	"github.com/urfave/cli/v2"
)

var (
	majorVersionReg = regexp.MustCompile(`^v?(\d+)\.\S+$`)
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:      "release",
		Aliases:   []string{"r"},
		Usage:     "List all releases of Flutter",
		Action:    handleRelease,
		UsageText: consts.APP_NAME + " release",
	})
}

func handleRelease(ctx *cli.Context) error {
	releases, err := utils.GetReleases()
	if err != nil {
		return err
	}

	majorVersionsMap := make(map[string][]model.Release, 0)
	for idx := range releases {
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
	sort.Strings(majorVersions)
	println()

	for majorIdx, majorVersion := range majorVersions {
		term.Success(fmt.Sprintf("[%s.x]:", majorVersion))
		count := 0
		printText := ""
		for _, release := range majorVersionsMap[majorVersion] {
			// Skip all pre-release or dev-release in outdated major version
			if majorIdx != len(majorVersions)-1 && strings.Contains(release.Version, "-") {
				continue
			}
			count++
			if count > 5 {
				break
			}
			printText += release.Version + "\n"
		}

		print(printText)
		if count > 5 {
			term.Warn(fmt.Sprintf("...and %d more", len(majorVersionsMap[majorVersion])-count))
		}
		println()
	}

	return nil
}
