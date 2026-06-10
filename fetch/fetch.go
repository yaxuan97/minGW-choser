package fetch

import (
	"encoding/json"
	"fmt"
	"mingw-chooser/match"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// GitHub API response types.
type githubRelease struct {
	Assets []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// buildNamePattern matches mingw-builds asset names.
// Examples:
//
//	x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z
//	i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z
var buildNamePattern = regexp.MustCompile(
	`^([a-z0-9_]+)-([0-9.]+)-release-([a-z0-9]+)-([a-z0-9]+)-([a-z0-9]+)-(.+)\.7z$`,
)

// Fetch retrieves the latest builds from the GitHub Releases API.
func Fetch(apiURL string) ([]match.Build, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch: create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "mingw-chooser")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch: API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("fetch: decode response: %w", err)
	}

	builds := buildListFromAssets(release.Assets)
	if len(builds) == 0 {
		return nil, fmt.Errorf("fetch: no .7z assets found in release")
	}

	return builds, nil
}

func buildListFromAssets(assets []githubAsset) []match.Build {
	var builds []match.Build
	for _, a := range assets {
		b, err := parseBuildFromAsset(a.Name, a.BrowserDownloadURL)
		if err != nil {
			continue // skip non-build assets
		}
		builds = append(builds, b)
	}
	return builds
}

func parseBuildFromAsset(name, url string) (match.Build, error) {
	// Lowercase for case-insensitive matching.
	lower := strings.ToLower(name)
	matches := buildNamePattern.FindStringSubmatch(lower)
	if matches == nil {
		return match.Build{}, fmt.Errorf("parse: %q does not match build naming pattern", name)
	}

	// For the Name field, keep the original casing from the asset.
	arch := matches[1]
	gcc := matches[2]
	thread := matches[3]
	exception := matches[4]
	crt := matches[5]
	// matches[6] is the rest (e.g., "rt_v12-rev0")

	return match.Build{
		Name:      name,
		URL:       url,
		Arch:      arch,
		Thread:    thread,
		Exception: exception,
		CRT:       crt,
		GCC:       gcc,
	}, nil
}
