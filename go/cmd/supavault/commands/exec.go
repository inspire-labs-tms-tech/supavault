package commands

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	errors2 "github.com/inspire-labs-tms-tech/supavault/pkg/helpers/auth/secrets/errors"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/supabase"
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/db"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var ExecCommand = &cli.Command{
	Name:      "exec",
	Usage:     "execute a command with arguments",
	ArgsUsage: "-- [arguments]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "no-login",
			Aliases:     []string{"n"},
			Usage:       "force the command to run even if not logged in",
			DefaultText: "false",
		},
		&cli.StringSliceFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "set environment variables (e.g., --env KEY=value)",
		},
	},
	Action: func(c *cli.Context) error {

		noLogin := c.Bool("no-login")

		if c.NArg() == 0 {
			return cli.Exit(color.RedString("no command provided to execute (example: supavault exec -- printenv)"), 1)
		}

		// build environment
		env := os.Environ()                      // start with system environment
		for _, e := range c.StringSlice("env") { // add any flag variables
			parts := strings.SplitN(e, "=", 2)
			if len(parts) != 2 {
				return cli.Exit(color.RedString("invalid environment variable format: %s (must be KEY=value)", e), 1)
			}
			key, value := parts[0], parts[1]

			// Validate the KEY
			if err := validateEnvVarName(key); err != nil {
				return cli.Exit(color.RedString("invalid environment variable name '%s': %v", key, err), 1)
			}

			env = append(env, fmt.Sprintf("%s=%s", key, value)) // append flag vars
		}

		client, clientErr := supabase.GetClient()
		if clientErr != nil {
			if clientErr.Error() == errors2.NewNotFoundError().Error() {
				if !noLogin {
					return cli.Exit(color.RedString("not logged in (use `supavault auth login` or re-run the command with the `--no-login` flag)"), 1)
				}
			} else {
				return cli.Exit(color.RedString(clientErr.Error()), 1)
			}
		}
		if client != nil {
			defer client.Close()

			var envs []db.Environment
			client.Authenticate()
			if err := client.Get("environments", envs); err != nil {
				return cli.Exit(color.RedString(err.Error()), 1)
			}
		}

		// Collect the command and arguments
		args := c.Args().Slice()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}
		return nil
	},
}

func validateEnvVarName(name string) error {
	// Regular expression to validate upper-case variable names
	varNameRegex := `^[A-Z_][A-Z0-9_]*$`
	matched, err := regexp.MatchString(varNameRegex, name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("must contain only upper-case letters, numbers, and underscores, and start with a letter or underscore")
	}
	return nil
}
