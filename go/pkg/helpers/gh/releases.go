package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/gh"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func ParseRelease(jsonData string) (*gh.Release, error) {
	var release gh.Release
	err := json.Unmarshal([]byte(jsonData), &release)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not parse release json: %s", err))
	}

	return &release, nil
}

func GetRelease(target *version.Version) (*gh.Release, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("https://api.github.com/repos/inspire-labs-tms-tech/supavault/releases/tags/%s", target.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var release gh.Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &release, nil
}

func GetLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/inspire-labs-tms-tech/supavault/releases/latest")
	if err != nil {
		return "", cli.Exit(color.RedString(err.Error()), 1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", cli.Exit(color.RedString("Error: received status code %d", resp.StatusCode), 1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", cli.Exit(color.RedString("Error reading response body: %s", err.Error()), 1)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", cli.Exit(color.RedString("Error parsing JSON: %s", err.Error()), 1)
	}

	return release.TagName, nil
}
