package main

import (
	"github.com/inspire-labs-tms-tech/supavault/cmd/supavault/commands"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	err := (&cli.App{
		Name:  "supavault",
		Usage: "A Supabase key-store",
		Commands: []*cli.Command{
			commands.AuthCommand,
		},
	}).Run(os.Args)
	if err != nil {
		panic(err)
	}
}
