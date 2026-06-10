package output

import (
	"bytes"
	"mingw-chooser/detect"
	"mingw-chooser/match"
	"strings"
	"testing"
)

func TestPrintText_GoldenPath(t *testing.T) {
	sys := detect.SystemInfo{OS: "windows", OSVersion: "Windows 10.0.22631", Arch: "x86_64"}
	r := match.MatchResult{
		Build: match.Build{
			Name: "x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z",
			URL:  "https://example.com/build.7z",
			Arch: "x86_64", Thread: "posix", Exception: "seh", CRT: "ucrt",
		},
		Explanation: []match.DimensionChoice{
			{Dimension: "arch", Choice: "x86_64", Reason: "your CPU is 64-bit", Manual: false},
			{Dimension: "thread", Choice: "posix", Reason: "best C++11 support", Manual: false},
			{Dimension: "exception", Choice: "seh", Reason: "optimal on x86_64", Manual: false},
			{Dimension: "crt", Choice: "ucrt", Reason: "modern runtime", Manual: false},
		},
		Score: 9,
	}

	var buf bytes.Buffer
	err := PrintResult(&buf, sys, r, FormatText, Flags{})
	if err != nil {
		t.Fatalf("PrintResult error: %v", err)
	}

	out := buf.String()
	required := []string{
		"Detected system:",
		"x86_64 (64-bit)",
		"Recommended build:",
		"x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z",
		"How to install:",
		"Why this build?",
	}
	for _, s := range required {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestPrintText_NoMatch(t *testing.T) {
	sys := detect.SystemInfo{OS: "linux", OSVersion: "Ubuntu 24.04", Arch: "aarch64"}
	r := match.MatchResult{}

	var buf bytes.Buffer
	err := PrintResult(&buf, sys, r, FormatText, Flags{})
	if err != nil {
		t.Fatalf("PrintResult error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No matching builds found") {
		t.Errorf("expected 'No matching builds found', got: %s", out)
	}
}

func TestPrintText_Wow64Warning(t *testing.T) {
	sys := detect.SystemInfo{OS: "windows", OSVersion: "Win10", Arch: "x86_64", IsWow64: true}
	r := match.MatchResult{
		Build: match.Build{Name: "b.7z", URL: "u", Arch: "x86_64", Thread: "posix", Exception: "seh", CRT: "ucrt"},
		Explanation: []match.DimensionChoice{
			{Dimension: "arch", Choice: "x86_64", Reason: "your CPU is 64-bit"},
		},
	}

	var buf bytes.Buffer
	PrintResult(&buf, sys, r, FormatText, Flags{})

	if !strings.Contains(buf.String(), "64-bit capable, running 32-bit process") {
		t.Errorf("output should warn about WoW64")
	}
}
