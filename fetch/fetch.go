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

// mingwBuildsPattern matches official mingw-builds asset names.
// Examples:
//
//	x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z
//	i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z
var mingwBuildsPattern = regexp.MustCompile(
	`^([a-z0-9_]+)-([0-9.]+)-release-([a-z0-9]+)-([a-z0-9]+)-([a-z0-9]+)-(.+)\.7z$`,
)

// winLibsPattern matches WinLibs asset names.
// Examples:
//
//	winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.7z
//	winlibs-i686-posix-dwarf-gcc-16.1.0-mingw-w64msvcrt-14.0.0-r2.zip
var winLibsPattern = regexp.MustCompile(
	`^winlibs-([a-z0-9_]+)-([a-z]+)-([a-z]+)-gcc-([0-9.]+)-mingw-w64([a-z]+)-([0-9.]+)-r([0-9]+)\.(7z|zip)$`,
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
		return nil, fmt.Errorf("fetch: no build assets found in release")
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
	// Try mingw-builds pattern first, then WinLibs.
	b, err := parseMingwBuilds(name, url)
	if err == nil {
		return b, nil
	}
	return parseWinLibs(name, url)
}

func parseMingwBuilds(name, url string) (match.Build, error) {
	lower := strings.ToLower(name)
	matches := mingwBuildsPattern.FindStringSubmatch(lower)
	if matches == nil {
		return match.Build{}, fmt.Errorf("parse: %q does not match mingw-builds pattern", name)
	}

	return match.Build{
		Name:      name,
		URL:       url,
		Arch:      matches[1],
		GCC:       matches[2],
		Thread:    matches[3],
		Exception: matches[4],
		CRT:       matches[5],
	}, nil
}

func parseWinLibs(name, url string) (match.Build, error) {
	lower := strings.ToLower(name)
	matches := winLibsPattern.FindStringSubmatch(lower)
	if matches == nil {
		return match.Build{}, fmt.Errorf("parse: %q does not match winlibs pattern", name)
	}

	// Groups: 1=arch, 2=thread, 3=exception, 4=gcc, 5=crt, 6=mingw_ver, 7=rev, 8=ext
	return match.Build{
		Name:      name,
		URL:       url,
		Arch:      matches[1],
		Thread:    matches[2],
		Exception: matches[3],
		GCC:       matches[4],
		CRT:       matches[5],
	}, nil
}
