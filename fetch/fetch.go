package fetch

import "mingw-chooser/match"

// Fetch retrieves the latest builds from the GitHub Releases API.
// Returns builds parsed from the asset list, or an error if unreachable.
func Fetch(apiURL string) ([]match.Build, error) {
	return nil, nil
}
