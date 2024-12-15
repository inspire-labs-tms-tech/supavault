package supavault

import (
	"github.com/inspire-labs-tms-tech/supavault/pkg"
	commands2 "github.com/inspire-labs-tms-tech/supavault/pkg/supavault/commands"
	"github.com/urfave/cli/v2"
)

var Supavault = &cli.App{
	Name:    "supavault",
	Usage:   "A Supabase key-store",
	Version: pkg.Version,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "verbose",
			Usage:       "print (top-level) verbose output; displays top-level panic/runtime errors which are suppressed by default",
			DefaultText: "false",
			Value:       false,
		},
	},
	Commands: []*cli.Command{
		commands2.AuthCommand,
		commands2.ExecCommand,
		commands2.ManageCommand,
		commands2.UpdateCommand,
	},
}
