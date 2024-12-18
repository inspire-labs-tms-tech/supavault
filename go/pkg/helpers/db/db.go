package db

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v2"
	"log"
)

type SchemaMigrationRecord struct {
	Version    string   // Corresponds to the version column in the table
	Statements []string // Corresponds to the statements column (if using text[])
	Name       string   // Corresponds to the name column
}

func GetInstalledVersion(pool *pgxpool.Pool) (*version.Version, error) {
	var installedVersion string
	if err := pool.QueryRow(context.Background(), "SELECT version FROM supavault.version_history ORDER BY at DESC LIMIT 1").Scan(&installedVersion); err != nil {
		if err.Error() == "no rows in result set" {
			_version, _ := version.NewVersion("0.0.0")
			return _version, nil
		}
		return nil, cli.Exit(color.RedString("Unable to get installed version: %v", err), 1)
	}
	if ver, err := version.NewVersion(installedVersion); err != nil {
		return nil, cli.Exit(color.RedString("Unable to parse installed version: %v", err), 1)
	} else {
		return ver, nil
	}
}

func Batch(pool *pgxpool.Pool, statements []string) error {
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

func GetExistingMigrations(pool *pgxpool.Pool) (map[string]SchemaMigrationRecord, error) {

	rows, err := pool.Query(context.Background(), `SELECT version, statements, name FROM supabase_migrations.schema_migrations`)
	if err != nil {
		return nil, cli.Exit(color.RedString("failed to query existing migrations: %v", err), 1)
	}
	defer rows.Close()

	migrations := make(map[string]SchemaMigrationRecord)
	for rows.Next() {
		var record SchemaMigrationRecord
		err := rows.Scan(&record.Version, &record.Statements, &record.Name)
		if err != nil {
			log.Fatalf("Failed to scan row: %v\n", err)
		}
		migrations[record.Version] = record
	}

	if rows.Err() != nil {
		return nil, cli.Exit(color.RedString("failed to query existing migrations: %v", rows.Err()), 1)
	}

	return migrations, nil
}
