package match

import (
	"encoding/json"
	"os"
	"testing"
)

// loadRules reads the real builds.json to keep tests in sync with data.
func loadRules(t *testing.T) Rules {
	t.Helper()
	data, err := os.ReadFile("../builds.json")
	if err != nil {
		t.Fatalf("failed to read builds.json: %v", err)
	}
	var cfg struct {
		Rules Rules `json:"rules"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse builds.json: %v", err)
	}
	return cfg.Rules
}

func makeBuilds() []Build {
	return []Build{
		{Name: "x86_64-14.2.0-release-posix-seh-ucrt", URL: "http://a", Arch: "x86_64", Thread: "posix", Exception: "seh", CRT: "ucrt", GCC: "14.2.0"},
		{Name: "x86_64-14.2.0-release-posix-seh-msvcrt", URL: "http://b", Arch: "x86_64", Thread: "posix", Exception: "seh", CRT: "msvcrt", GCC: "14.2.0"},
		{Name: "x86_64-14.2.0-release-win32-seh-ucrt", URL: "http://c", Arch: "x86_64", Thread: "win32", Exception: "seh", CRT: "ucrt", GCC: "14.2.0"},
		{Name: "i686-14.2.0-release-posix-dwarf-ucrt", URL: "http://d", Arch: "i686", Thread: "posix", Exception: "dwarf", CRT: "ucrt", GCC: "14.2.0"},
		{Name: "i686-14.1.0-release-posix-dwarf-ucrt", URL: "http://e", Arch: "i686", Thread: "posix", Exception: "dwarf", CRT: "ucrt", GCC: "14.1.0"},
	}
}

func TestMatch_PrefersPosixSehUcrtOnX86_64(t *testing.T) {
	builds := makeBuilds()
	rules := loadRules(t)
	result := Match("x86_64", builds, rules, Overrides{})

	if result.Build.Thread != "posix" {
		t.Errorf("expected posix thread, got %s", result.Build.Thread)
	}
	if result.Build.Exception != "seh" {
		t.Errorf("expected seh exception, got %s", result.Build.Exception)
	}
	if result.Build.CRT != "ucrt" {
		t.Errorf("expected ucrt CRT, got %s", result.Build.CRT)
	}
	if result.Build.Name != "x86_64-14.2.0-release-posix-seh-ucrt" {
		t.Errorf("expected best build, got %s", result.Build.Name)
	}
}

func TestMatch_FiltersByArch(t *testing.T) {
	builds := makeBuilds()
	rules := loadRules(t)
	result := Match("i686", builds, rules, Overrides{})

	if result.Build.Arch != "i686" {
		t.Errorf("expected i686 arch, got %s", result.Build.Arch)
	}
	for _, alt := range result.Alternatives {
		if alt.Arch != "i686" {
			t.Errorf("alternative has wrong arch: %s", alt.Arch)
		}
	}
}

func TestMatch_GccTiebreaker(t *testing.T) {
	builds := []Build{
		{Name: "old", URL: "u", Arch: "i686", Thread: "posix", Exception: "dwarf", CRT: "ucrt", GCC: "13.2.0"},
		{Name: "new", URL: "u", Arch: "i686", Thread: "posix", Exception: "dwarf", CRT: "ucrt", GCC: "14.2.0"},
	}
	rules := loadRules(t)
	result := Match("i686", builds, rules, Overrides{})

	if result.Build.Name != "new" {
		t.Errorf("expected newer GCC to win tiebreaker, got %s", result.Build.Name)
	}
}

func TestMatch_UserOverrideForcesDimension(t *testing.T) {
	builds := makeBuilds()
	rules := loadRules(t)
	result := Match("x86_64", builds, rules, Overrides{Thread: "win32"})

	if result.Build.Thread != "win32" {
		t.Errorf("expected win32 override, got %s", result.Build.Thread)
	}
	for _, exp := range result.Explanation {
		if exp.Dimension == "thread" && !exp.Manual {
			t.Error("thread dimension should be marked Manual=true when overridden")
		}
	}
}

func TestMatch_NoMatchingBuilds(t *testing.T) {
	builds := makeBuilds()
	rules := loadRules(t)
	result := Match("aarch64", builds, rules, Overrides{})

	if result.Build.Name != "" {
		t.Errorf("expected empty build when no matches, got %s", result.Build.Name)
	}
}

func TestMatch_ExplanationHasAllDimensions(t *testing.T) {
	builds := makeBuilds()
	rules := loadRules(t)
	result := Match("x86_64", builds, rules, Overrides{})

	dims := make(map[string]bool)
	for _, exp := range result.Explanation {
		dims[exp.Dimension] = true
	}
	required := []string{"arch", "thread", "exception", "crt"}
	for _, r := range required {
		if !dims[r] {
			t.Errorf("missing explanation for dimension %s", r)
		}
	}
}

func TestMatch_ScoreIsCorrect(t *testing.T) {
	builds := makeBuilds()
	rules := loadRules(t)
	result := Match("x86_64", builds, rules, Overrides{})

	// x86_64 posix seh ucrt = 1st in all 3 scored dims = 3+3+3 = 9
	if result.Score != 9 {
		t.Errorf("expected score 9 for perfect match, got %d", result.Score)
	}
}
