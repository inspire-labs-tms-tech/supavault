package main

import (
	"github.com/inspire-labs-tms-tech/supavault/cmd/supavault/commands"
	"github.com/inspire-labs-tms-tech/supavault/pkg"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	err := (&cli.App{
		Name:    "supavault",
		Usage:   "A Supabase key-store",
		Version: pkg.Version,
		Commands: []*cli.Command{
			commands.AuthCommand,
			commands.ExecCommand,
			commands.ManageCommand,
			commands.UpdateCommand,
		},
	}).Run(os.Args)
	if err != nil {
		panic(err)
	}
}
