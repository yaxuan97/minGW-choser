package output

import (
	"bytes"
	"encoding/json"
	"mingw-chooser/detect"
	"mingw-chooser/match"
	"strings"
	"testing"
)

func TestPrintJSON_ValidOutput(t *testing.T) {
	sys := detect.SystemInfo{OS: "windows", OSVersion: "10.0.22631", Arch: "x86_64"}
	r := match.MatchResult{
		Build: match.Build{
			Name: "b.7z", URL: "http://u",
			Arch: "x86_64", Thread: "posix", Exception: "seh", CRT: "ucrt",
		},
		Alternatives: []match.Build{
			{Name: "alt.7z", URL: "http://u2", Arch: "x86_64", Thread: "win32", Exception: "seh", CRT: "ucrt"},
		},
		Explanation: []match.DimensionChoice{
			{Dimension: "arch", Choice: "x86_64", Reason: "your CPU is 64-bit"},
		},
		Score: 9,
	}

	var buf bytes.Buffer
	err := PrintResult(&buf, sys, r, FormatJSON, Flags{})
	if err != nil {
		t.Fatalf("PrintResult JSON error: %v", err)
	}

	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if out.System.OS != "windows" {
		t.Errorf("expected OS windows, got %s", out.System.OS)
	}
	if out.Recommended == nil || out.Recommended.Name != "b.7z" {
		t.Error("expected recommended build")
	}
	if len(out.Alternatives) != 1 {
		t.Errorf("expected 1 alternative, got %d", len(out.Alternatives))
	}
}

func TestPrintJSON_WoW64Warning(t *testing.T) {
	sys := detect.SystemInfo{OS: "windows", OSVersion: "10", Arch: "x86_64", IsWow64: true}
	r := match.MatchResult{
		Build: match.Build{Name: "b.7z", URL: "u", Arch: "x86_64", Thread: "posix", Exception: "seh", CRT: "ucrt"},
	}

	var buf bytes.Buffer
	PrintResult(&buf, sys, r, FormatJSON, Flags{})

	var out jsonOutput
	json.Unmarshal(buf.Bytes(), &out)
	if !strings.Contains(out.Warning, "32-bit process") {
		t.Errorf("expected WoW64 warning, got: %s", out.Warning)
	}
}

func TestPrintJSON_NoMatch(t *testing.T) {
	sys := detect.SystemInfo{OS: "linux", OSVersion: "24.04", Arch: "aarch64"}
	r := match.MatchResult{}

	var buf bytes.Buffer
	err := PrintResult(&buf, sys, r, FormatJSON, Flags{})
	if err != nil {
		t.Fatalf("PrintResult JSON error: %v", err)
	}

	var out jsonOutput
	json.Unmarshal(buf.Bytes(), &out)
	if out.Recommended != nil {
		t.Error("expected nil recommended for no match")
	}
}
