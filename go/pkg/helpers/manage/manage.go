package manage

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/inspire-labs-tms-tech/supavault/pkg"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/db"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/gh"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func SetupSupavault(pool *pgxpool.Pool, _targetVersion string, verbose bool) (*version.Version, error) {

	installedVersion, err := db.GetInstalledVersion(pool)
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
	migrations, err := db.GetExistingMigrations(pool)
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
			if err := db.Batch(pool, statements); err != nil {
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
