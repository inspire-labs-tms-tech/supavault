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
	"runtime"
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
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Usage:       "verbose output (useful for debugging)",
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
		verbose := c.Bool("verbose")

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
			client.Authenticate()

			user, _ := client.GetUser()
			if verbose {
				color.Blue("user: %v", user.ID)
			}

			var projects []db.Project
			if err := client.Get("projects", &projects); err != nil {
				return cli.Exit(color.RedString(err.Error()), 1)
			} else if projects == nil {
				return cli.Exit(color.RedString("no projects found"), 1)
			} else if len(projects) == 0 {
				return cli.Exit(color.RedString("no projects found"), 1)
			} else if len(projects) != 1 {
				return cli.Exit(color.RedString("multiple projects found"), 1)
			}
			project := projects[0]
			if verbose {
				color.Blue("project: %s", project.ID.String())
			}

			var environments []db.Environment
			if err := client.Get("environments", &environments); err != nil {
				return cli.Exit(color.RedString(err.Error()), 1)
			} else if environments == nil {
				return cli.Exit(color.RedString("no environments found"), 1)
			} else if len(environments) != 1 {
				return cli.Exit(color.RedString("multiple environments found"), 1)
			}
			environment := environments[0]
			if verbose {
				color.Blue("environment: %s", environment.ID.String())
			}

			var vars []db.Variable
			if err := client.Get("variables", &vars); err != nil {
				return cli.Exit(color.RedString(err.Error()), 1)
			}
			varMap := make(map[string]string)
			for _, v := range vars {
				varMap[v.ID] = v.Default
			}
			if verbose {
				color.Blue("variables: %v", len(vars))
			}

			var envVars []db.EnvironmentVariable
			if err := client.Get("environment_variables", &envVars); err != nil {
				return cli.Exit(color.RedString(err.Error()), 1)
			}
			if verbose {
				color.Blue("environment variables: %v", len(envVars))
			}

			for _, envVar := range envVars {
				if envVar.Value == "" { // use default if none provided
					env = append(env, fmt.Sprintf("%s=%s", envVar.VariableID, varMap[envVar.VariableID]))
				} else {
					env = append(env, fmt.Sprintf("%s=%s", envVar.VariableID, envVar.Value))
				}
			}

		}

		// prepare the command
		args := c.Args().Slice()
		commandString := strings.Join(args, " ")
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", commandString)
		} else {
			cmd = exec.Command("/bin/sh", "-c", commandString)
		}
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if verbose {
			color.Blue("command: %s", commandString)
			color.Blue("===== start command output =====")
		}
		err := cmd.Run()
		if verbose {
			color.Blue("====== end command output ======")
		}

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
