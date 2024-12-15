package commands

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	commands "github.com/inspire-labs-tms-tech/supavault/pkg/helpers"
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
		Usage:       "force the command to continue if certain (overcome-able) errors occur",
		DefaultText: "false",
		Value:       false,
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
						return cli.Exit(color.RedString("invalid password (enter a valid password or specify one with the --%s flag)", PASS), 1)
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
					return cli.Exit(color.RedString("Unable to connect to database: %v", err), 1)
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
				if err := setupSupavaultSchemaIfNotExists(pool); err != nil {
					color.Red("unable to ensure supavault schema is configured")
					return cli.Exit(color.RedString(err.Error()), 1)
				}

				// ensure supabase_migrations schema is setup and exists
				if verbose {
					color.Blue("ensuring supabase_migrations schema exists...")
				}
				if err := setupSupabaseMigrationsIfNotExists(pool); err != nil {
					color.Red("unable to ensure supabase_migrations schema is configured")
					return cli.Exit(color.RedString(err.Error()), 1)
				}

				color.Green("install complete!")

				return nil
			},
		},
	},
}

func batch(pool *pgxpool.Pool, statements []string) error {
	txn, err := pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf(" -> unable to start transaction\n   -> %w", err)
	}

	for _, stmt := range statements {
		_, err := txn.Exec(context.Background(), stmt)
		if err != nil {
			_ = txn.Rollback(context.Background())
			return fmt.Errorf(" -> failed to execute statement \"%s\"\n   -> %w", stmt, err)
		}
	}

	if err := txn.Commit(context.Background()); err != nil {
		return fmt.Errorf(" -> failed to commit transaction\n   -> %w", err)
	}

	return nil
}

func setupSupavaultSchemaIfNotExists(pool *pgxpool.Pool) error {
	statements := []string{
		"CREATE SCHEMA IF NOT EXISTS supavault",
		"CREATE TABLE IF NOT EXISTS supavault.version_history ()",
		"ALTER TABLE supavault.version_history ADD COLUMN IF NOT EXISTS version text NOT NULL PRIMARY KEY",
		"ALTER TABLE supavault.version_history ADD COLUMN IF NOT EXISTS at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP",
		"ALTER TABLE supavault.version_history DROP CONSTRAINT IF EXISTS valid_semver",
		"ALTER TABLE supavault.version_history ADD CONSTRAINT valid_semver CHECK (version ~ '^(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(-(0|[1-9A-Za-z-][0-9A-Za-z-]*)(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$')",
	}

	return batch(pool, statements)
}

func setupSupabaseMigrationsIfNotExists(pool *pgxpool.Pool) error {
	statements := []string{
		"CREATE SCHEMA IF NOT EXISTS supabase_migrations",
		"CREATE TABLE IF NOT EXISTS supabase_migrations.schema_migrations ()",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS version text NOT NULL PRIMARY KEY",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS statements text[]",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS name text",
		"CREATE TABLE IF NOT EXISTS supabase_migrations.seed_files ()",
		"ALTER TABLE supabase_migrations.seed_files ADD COLUMN IF NOT EXISTS path text NOT NULL PRIMARY KEY",
		"ALTER TABLE supabase_migrations.seed_files ADD COLUMN IF NOT EXISTS hash text NOT NULL",
	}

	return batch(pool, statements)
}
