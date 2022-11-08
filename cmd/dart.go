package cmd

import "github.com/urfave/cli/v2"

func init() {
	cmds = append(cmds, &cli.Command{
		Name:  "dart",
		Usage: "Proxy dart commands",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	})
}
