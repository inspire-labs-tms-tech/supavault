package main

import (
	"github.com/fatih/color"
	"github.com/inspire-labs-tms-tech/supavault/pkg/supavault"
	"os"
)

func main() {
	if err := supavault.Supavault.Run(os.Args); err != nil {
		if len(os.Args) > 0 && os.Args[0] == "--verbose" {
			color.Red(err.Error())
		}
		os.Exit(1)
	}
	os.Exit(0)
}
