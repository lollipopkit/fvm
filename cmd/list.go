package cmd

import "github.com/urfave/cli/v2"

func init() {
	cmds = append(cmds, &cli.Command{
		Name:  "list",
		Usage: "List all installed versions of Flutter",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	})
}
