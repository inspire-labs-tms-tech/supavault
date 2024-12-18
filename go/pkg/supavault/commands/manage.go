package commands

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/inspire-labs-tms-tech/supavault/pkg"
	commands "github.com/inspire-labs-tms-tech/supavault/pkg/helpers"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/gh"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
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

type SchemaMigrationRecord struct {
	Version    string   // Corresponds to the version column in the table
	Statements []string // Corresponds to the statements column (if using text[])
	Name       string   // Corresponds to the name column
}

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

				// download and apply migrations
				if newVersion, err := setupSupavault(pool, c.String("version"), verbose); err != nil {
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

func getInstalledVersion(pool *pgxpool.Pool) (*version.Version, error) {
	var installedVersion string
	if err := pool.QueryRow(context.Background(), "SELECT version FROM supavault.version_history ORDER BY at DESC LIMIT 1").Scan(&installedVersion); err != nil {
		return nil, cli.Exit(color.RedString("Unable to get installed version: %v", err), 1)
	}
	if ver, err := version.NewVersion(installedVersion); err != nil {
		return nil, cli.Exit(color.RedString("Unable to parse installed version: %v", err), 1)
	} else {
		return ver, nil
	}
}

func getTargetVersion(_targetVersion string) (*version.Version, error) {
	if _targetVersion == "" || _targetVersion == "latest" {
		latestVersion, latestVersionError := gh.GetLatestVersion()
		if latestVersionError != nil {
			return nil, cli.Exit(color.RedString("unable to determine the latest version: %s", latestVersionError.Error()), 1)
		}
		_targetVersion = latestVersion
	}
	if targetVersion, err := version.NewVersion(_targetVersion); err != nil {
		return nil, cli.Exit(color.RedString(err.Error()), 1)
	} else if targetVersion == nil {
		return nil, cli.Exit(color.RedString("invalid version: %s", _targetVersion), 1)
	} else {
		return targetVersion, nil
	}
}

func compareVersions(installedVersion *version.Version, targetVersion *version.Version) error {
	if targetVersion.GreaterThan(installedVersion) {
		return nil
	} else if targetVersion.Equal(installedVersion) {
		return cli.Exit(color.YellowString("current version (%s) is target version (%s)", pkg.Version, targetVersion.String()), 0)
	} else {
		return cli.Exit(color.RedString("target version (%s) is older than installed version (%s)", targetVersion.String(), pkg.Version), 1)
	}
}

func setupSupavault(pool *pgxpool.Pool, _targetVersion string, verbose bool) (*version.Version, error) {

	installedVersion, err := getInstalledVersion(pool)
	if err != nil {
		return nil, err
	}

	targetVersion, err := getTargetVersion(_targetVersion)
	if err != nil {
		return nil, err
	}

	if err := compareVersions(installedVersion, targetVersion); err != nil {
		return nil, err
	}

	downloadURL := fmt.Sprintf("https://github.com/inspire-labs-tms-tech/supavault/releases/download/%s/migrations.zip", targetVersion.String())
	if zipData, err := downloadZip(downloadURL); err != nil {
		return nil, cli.Exit(color.RedString("failed to download ZIP: %w", err), 1)
	} else if err := processZip(zipData, pool, verbose); err != nil {
		return nil, err
	}

	return targetVersion, nil
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

func downloadZip(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, cli.Exit(color.RedString("HTTP GET failed: %w", err), 1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, cli.Exit(color.RedString("HTTP GET failed: %s", resp.Status), 1)
	}

	return ioutil.ReadAll(resp.Body)
}

func readZipFile(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", cli.Exit(color.RedString("failed to open zip file: %v", err), 1)
	}
	defer rc.Close()

	content, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", cli.Exit(color.RedString("failed to read zip file: %v", err), 1)
	}
	return string(content), nil
}

func getExistingMigrations(pool *pgxpool.Pool) (map[string]SchemaMigrationRecord, error) {

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

func processZip(zipData []byte, pool *pgxpool.Pool, verbose bool) error {

	// read the ZIP archive
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return cli.Exit(color.RedString("failed to open ZIP reader: %w", err), 1)
	}

	files := make(map[string]string)
	fileNames := []string{}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		fileNames = append(fileNames, file.Name)
		content, err := readZipFile(file)
		if err != nil {
			return cli.Exit(color.RedString("failed to read file '%s': %w", file.Name, err), 1)
		}
		files[file.Name] = content
	}
	sort.Strings(fileNames) // sort files by name

	// get already-applied migrations
	migrations, err := getExistingMigrations(pool)
	if err != nil {
		return err
	}

	// Process sorted migration files
	for _, fileName := range fileNames {
		pathParts := strings.Split(fileName, "/")
		parts := strings.Split(pathParts[len(pathParts)-1], "_")
		v := parts[0]
		description := strings.TrimSuffix(parts[1], ".sql")
		sql := files[fileName]

		if record, exists := migrations[v]; exists {
			if verbose {
				color.Blue("version %s already applied...", record.Version)
			}
		} else {
			var statements []string
			statements = append(statements, sql)
			statements = append(statements, fmt.Sprintf(
				"INSERT INTO supabase_migrations.schema_migrations (version, statements, name) VALUES ('%s', ARRAY[%s], '%s')",
				v,
				formatStatementsForSQL(statements),
				description,
			))
			if err := batch(pool, statements); err != nil {
				return cli.Exit(color.RedString("failed to execute migration '%s': \n%s", v, err.Error()), 1)
			}
		}
	}
	return nil
}

// Helper function to format statements for SQL
func formatStatementsForSQL(statements []string) string {
	var formatted []string
	for _, stmt := range statements {
		// Escape single quotes
		stmt = strings.ReplaceAll(stmt, "'", "''")

		// Wrap the statement in $escaped$...$escaped$ to handle potential dollar-quoted strings
		formatted = append(formatted, fmt.Sprintf("$escaped$%s$escaped$", stmt))
	}
	return strings.Join(formatted, ", ")
}
