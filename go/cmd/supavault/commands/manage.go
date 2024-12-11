package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var ManageCommand = &cli.Command{
	Name:        "manage",
	Description: "manage a supavault instance of Supabase",
	Subcommands: []*cli.Command{
		{
			Name:        "install",
			Description: "install supavault in a Supabase instance",
			Usage:       "install supavault in a Supabase instance",
			Action: func(c *cli.Context) error {

				// dsn := "postgres://username:password@localhost:5432/database_name"
				conn, err := pgxpool.New(context.Background(), dsn)
				if err != nil {
					return cli.Exit(color.RedString("Unable to connect to database: %v", err), 1)
				}
				defer conn.Close()

				// ensure supabase_migrations schema is setup and exists
				if err := setupSupabaseMigrationsIfNotExists(conn); err != nil {
					return cli.Exit(color.RedString(err.Error()), 1)
				}

				return nil
			},
		},
	},
}

func setupSupabaseMigrationsIfNotExists(conn *pgxpool.Pool) error {
	statements := []string{
		"CREATE SCHEMA IF NOT EXISTS supabase_migrations",
		"CREATE TABLE IF NOT EXISTS supabase_migrations.schema_migrations (version text NOT NULL PRIMARY KEY)",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS statements text[]",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS name text",
		"CREATE TABLE IF NOT EXISTS supabase_migrations.seed_files (path text NOT NULL PRIMARY KEY, hash text NOT NULL)",
	}

	for _, stmt := range statements {
		_, err := conn.Exec(context.Background(), stmt)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to verify supabase_migrations schema (statement: '%v'): %v", stmt, err))
		}
	}
	return nil
}
