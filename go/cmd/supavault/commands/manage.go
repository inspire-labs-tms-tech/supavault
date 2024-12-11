package commands

import "github.com/urfave/cli/v2"

var ManageCommand = &cli.Command{
	Name:        "manage",
	Description: "manage a supavault instance of Supabase",
	Subcommands: []*cli.Command{
		{
			Name:        "install",
			Description: "install supavault in a Supabase instance",
			Usage:       "install supavault in a Supabase instance",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	},
}
