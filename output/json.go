package output

import (
	"encoding/json"
	"io"
	"mingw-chooser/detect"
	"mingw-chooser/match"
)

type jsonOutput struct {
	System       systemJSON            `json:"system"`
	Recommended  *buildJSON            `json:"recommended"`
	Alternatives []buildJSON           `json:"alternatives,omitempty"`
	Explanation  []match.DimensionChoice `json:"explanation"`
	Warning      string                `json:"warning,omitempty"`
}

type systemJSON struct {
	OS        string `json:"os"`
	OSVersion string `json:"os_version"`
	Arch      string `json:"arch"`
}

type buildJSON struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func printJSON(w io.Writer, sys detect.SystemInfo, r match.MatchResult, flags Flags) error {
	out := jsonOutput{
		System: systemJSON{
			OS:        sys.OS,
			OSVersion: sys.OSVersion,
			Arch:      sys.Arch,
		},
		Explanation: r.Explanation,
	}

	if r.Build.Name != "" {
		out.Recommended = &buildJSON{Name: r.Build.Name, URL: r.Build.URL}
	}

	for _, alt := range r.Alternatives {
		out.Alternatives = append(out.Alternatives, buildJSON{Name: alt.Name, URL: alt.URL})
	}

	if sys.IsWow64 {
		out.Warning = "32-bit process detected on 64-bit system; recommending 64-bit build"
	}

	if flags.Offline {
		if out.Warning != "" {
			out.Warning += "; "
		}
		out.Warning += "offline mode: using embedded build snapshot, may not be the latest version"
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
