package commands

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/inspire-labs-tms-tech/supavault/pkg"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/auth/sudo"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/gh"
	gh2 "github.com/inspire-labs-tms-tech/supavault/pkg/models/gh"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var osName = runtime.GOOS
var arch = runtime.GOARCH

var UpdateCommand = &cli.Command{
	Name:        "update",
	Usage:       "update the supavault CLI installed on the local machine",
	Description: "update the supavault CLI installed on the local machine",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "version",
			Aliases:     []string{},
			Usage:       "the version to install (`latest` or a semver)",
			DefaultText: "latest",
		},
	},
	Action: func(c *cli.Context) error {

		if !sudo.IsElevated() {
			return cli.Exit(color.RedString("command must be run as root user"), 1)
		}

		_targetVersion := c.String("version")
		if _targetVersion == "" || _targetVersion == "latest" {
			latestVersion, latestVersionError := gh.GetLatestVersion()
			if latestVersionError != nil {
				return cli.Exit(color.RedString("Unable to determine the latest version: %s", latestVersionError.Error()), 1)
			}
			_targetVersion = latestVersion
		}
		installedVersion, _ := version.NewVersion(pkg.Version)
		targetVersion, vErr := version.NewVersion(_targetVersion)
		if vErr != nil {
			return cli.Exit(color.RedString(vErr.Error()), 1)
		}
		if targetVersion == nil {
			return cli.Exit(color.RedString("invalid version"), 1)
		}

		if targetVersion.GreaterThan(installedVersion) {
			doUpdate(targetVersion)
		} else if targetVersion.Equal(installedVersion) {
			color.Yellow("Current version is target version")
		} else {
			cli.Exit(color.RedString("Error: target version is older than installed version"), 1)
		}

		return nil
	},
}

func doUpdate(targetVersion *version.Version) {
	if targetVersion == nil {
		panic("target version is nil")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1) // Buffered channel to collect errors

	installDir := filepath.Join("/", "usr", "local", "bin", "supavault")
	if osName == "windows" {
		installDir = filepath.Join(os.Getenv("APPDATA"), "Local", "Programs", "supavault")
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = color.BlueString(" Updating to %s...", targetVersion.String())

	// Start the spinner
	s.Start()

	// Goroutine for the update process
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Fetch the release
		release, releaseErr := gh.GetRelease(targetVersion)
		if releaseErr != nil {
			errChan <- fmt.Errorf("failed to fetch release: %w", releaseErr)
			return
		}

		// Locate the target file in release assets
		targetFile := fmt.Sprintf("supavault_%s_%s_%s", osName, arch, targetVersion.String())
		var releaseFile *gh2.Asset
		for _, asset := range release.Assets {
			if asset.Name == targetFile {
				releaseFile = &asset
				break
			}
		}
		if releaseFile == nil {
			errChan <- fmt.Errorf("unable to find the release file: %s", targetFile)
			return
		}

		// Download the file
		err := downloadFile(releaseFile.BrowserDownloadURL, installDir)
		if err != nil {
			errChan <- fmt.Errorf("failed to download file: %w", err)
			return
		}

		// Signal successful completion
		errChan <- nil
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Stop the spinner
	s.Stop()

	// Handle errors
	close(errChan) // Close the channel to signal completion
	for err := range errChan {
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}

	// Success message
	color.Green("Update completed and file downloaded successfully!")
}

func downloadFile(url, fp string) error {
	out, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Step 2: Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Step 3: Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Step 4: Write the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if err := os.Chmod(fp, 0755); err != nil {
		return fmt.Errorf("failed to make file executable: %w", err)
	}

	return nil
}
