package cmd

import "github.com/urfave/cli/v2"

func init() {
	cmds = append(cmds, &cli.Command{
		Name:  "global",
		Usage: "Manage global version of Flutter",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	})
}
