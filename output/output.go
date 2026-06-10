package output

import (
	"io"
	"mingw-chooser/detect"
	"mingw-chooser/match"
)

// Format selects output rendering mode.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Flags carries user-supplied overrides for display purposes.
type Flags struct {
	Arch      string
	Thread    string
	Exception string
	CRT       string
	Offline   bool
	JSON      bool
	List      bool
}

// PrintResult writes the match result to w in the requested format.
func PrintResult(w io.Writer, sys detect.SystemInfo, r match.MatchResult, f Format, flags Flags) error {
	switch f {
	case FormatJSON:
		return printJSON(w, sys, r, flags)
	default:
		return printText(w, sys, r, flags)
	}
}
