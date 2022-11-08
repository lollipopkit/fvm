package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/LollipopKit/gofvm/consts"
	"github.com/LollipopKit/gofvm/model"
	"github.com/LollipopKit/gofvm/term"
	"github.com/LollipopKit/gofvm/utils"
	"github.com/urfave/cli/v2"
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
		v := strings.Split(releases[idx].Version, ".")
		if len(v) < 3 {
			continue
		}
		majorVersion := v[0]
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

	term.Cyan("\nFlutter releases:")

	for _, majorVersion := range majorVersions {
		term.Green(fmt.Sprintf("[%s.x]:", majorVersion))
		count := 0
		printText := ""
		for idx := range majorVersionsMap[majorVersion] {
			count++
			if count > 5 {
				break
			}
			printText += majorVersionsMap[majorVersion][idx].Version + "\n"
		}

		print(printText)
		if count > 5 {
			term.Yellow(fmt.Sprintf("...and %d more", len(majorVersionsMap[majorVersion])-count))
		}
		println()
	}

	return nil
}
