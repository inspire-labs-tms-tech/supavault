package commands

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	commands "github.com/inspire-labs-tms-tech/supavault/pkg/helpers"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/db"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/manage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v2"
	"strings"
)

const CONN_STR = "conn-str"

const USER = "db-user"
const USER_DEFAULT = "postgres"

const PASS = "db-pass"

const NAME = "db-name"
const NAME_DEFAULT = "postgres"

const HOST = "db-host"
const HOST_DEFAULT = "127.0.0.1"

const PORT = "db-port"
const PORT_DEFAULT = 5432

var flags = &[]cli.Flag{
	&cli.StringFlag{
		Name:    CONN_STR,
		Aliases: []string{},
		Usage:   "connection string of the database (override other connection values if user; e.g., postgres://user:pass@host:port/db)",
	},
	&cli.StringFlag{
		Name:        USER,
		Aliases:     []string{},
		Usage:       "the name of the database user to authenticate with",
		DefaultText: USER_DEFAULT,
	},
	&cli.StringFlag{
		Name:    PASS,
		Aliases: []string{},
		Usage:   "the password of the database user to authenticate with",
	},
	&cli.StringFlag{
		Name:        NAME,
		Aliases:     []string{},
		Usage:       "the name of the database to authenticate into",
		DefaultText: NAME_DEFAULT,
	},
	&cli.StringFlag{
		Name:    HOST,
		Aliases: []string{},
		Usage:   "the host of the database to authenticate into",
	},
	&cli.IntFlag{
		Name:    PORT,
		Aliases: []string{},
		Usage:   "the port of the database to authenticate into",
	},
	&cli.BoolFlag{
		Name:        "verbose",
		Usage:       "print verbose output",
		DefaultText: "false",
		Value:       false,
	},
	&cli.BoolFlag{
		Name:        "force",
		Aliases:     []string{"update"},
		Usage:       "force the command to continue if certain (overcome-able) errors occur",
		DefaultText: "false",
		Value:       false,
	},
	&cli.StringFlag{
		Name:        "version",
		Aliases:     []string{},
		Usage:       "the version to install (`latest` or a semver)",
		DefaultText: "latest",
	},
}

var ManageCommand = &cli.Command{
	Name:        "manage",
	Description: "manage a supavault instance of Supabase",
	Usage:       "manage a supavault instance of Supabase",
	Subcommands: []*cli.Command{
		{
			Name:        "install",
			Description: "install supavault in a Supabase instance",
			Usage:       "install supavault in a Supabase instance",
			Flags:       *flags,
			Action: func(c *cli.Context) error {

				verbose := c.Bool("verbose")
				force := c.Bool("force")

				pool, connErr := getPool(c)
				if connErr != nil {
					return connErr
				}
				defer pool.Close()

				// Query to check if the schema exists
				if verbose {
					color.Blue("checking if supavault is installed...")
				}
				var exists bool
				if err := pool.QueryRow(
					context.Background(),
					"SELECT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'supavault')",
				).Scan(&exists); err != nil {
					return cli.Exit(color.RedString("Unable to check if supavault schema exists: %v", err), 1)
				}
				if exists {
					if !force {
						return cli.Exit(color.RedString("supavault already installed (use the update command to install a newer version or use the --force flag to continue)"), 1)
					} else if verbose {
						color.Yellow("continuing with forced install...")
					}
				}

				// ensure own schema is setup and exists
				if verbose {
					color.Blue("ensuring supavault schema exists...")
				}
				if err := db.SetupSupavaultSchemaIfNotExists(pool); err != nil {
					color.Red("unable to ensure supavault schema is configured")
					return cli.Exit(color.RedString(err.Error()), 1)
				}

				// ensure supabase_migrations schema is setup and exists
				if verbose {
					color.Blue("ensuring supabase_migrations schema exists...")
				}
				if err := db.SetupSupabaseMigrationsSchemaIfNotExists(pool); err != nil {
					color.Red("unable to ensure supabase_migrations schema is configured")
					return cli.Exit(color.RedString(err.Error()), 1)
				}

				// download and apply migrations
				if newVersion, err := manage.SetupSupavault(pool, c.String("version"), verbose); err != nil {
					color.Red("unable to install supavault")
					return cli.Exit(color.RedString(err.Error()), 1)
				} else if verbose {
					color.Green("installed supavault %v", newVersion.String())
				}

				color.Green("install complete!")

				return nil
			},
		},
	},
}

func getPool(c *cli.Context) (*pgxpool.Pool, error) {

	verbose := c.Bool("verbose")
	connStr := strings.TrimSpace(c.String(CONN_STR))

	user := c.String(USER)
	if user == "" {
		user = USER_DEFAULT
	}

	pass := strings.TrimSpace(c.String(PASS))
	if pass == "" && connStr == "" {
		prompted, _ := commands.PromptRaw("Password")
		pass = prompted
		pass = strings.TrimSpace(pass)
		if pass == "" {
			return nil, cli.Exit(color.RedString("invalid password (enter a valid password or specify one with the --%s flag)", PASS), 1)
		}
	}

	name := c.String(NAME)
	if name == "" {
		name = NAME_DEFAULT
	}

	host := c.String(HOST)
	if host == "" {
		host = HOST_DEFAULT
	}

	port := c.Int(PORT)
	if port == 0 {
		port = PORT_DEFAULT
	}

	if connStr == "" {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", name, pass, host, port, name)
	}

	if verbose {
		color.Blue("authenticating with: postgres://%s:%s@%s:%d/%s", name, strings.Repeat("*", len(pass)), host, port, name)
	}

	// create connection pool
	if verbose {
		color.Blue("creating connection pool...")
	}
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, cli.Exit(color.RedString("Unable to connect to database: %v", err), 1)
	}
	return pool, nil
}
