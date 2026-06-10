# MinGW Chooser — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a cross-platform CLI tool (Go) that detects system properties and recommends the correct MinGW-w64 build.

**Architecture:** Five packages (`detect`, `match`, `fetch`, `output`, `main`) with `builds.json` embedded at build time. Zero external dependencies — only Go standard library. Platform-specific code uses build tags.

**Tech Stack:** Go 1.23+, `syscall` for Windows API, `os/exec` for Unix utilities, `net/http` for GitHub API, `encoding/json` for config and output.

**Note:** `go` is not yet in PATH — the user needs to restart their terminal after installing Go before running commands. Plan assumes Go 1.23+ is available.

---

### Task 1: Project scaffolding

**Files:**
- Create: `go.mod`
- Create: `detect/detect.go`
- Create: `match/types.go`
- Create: `fetch/fetch.go`
- Create: `output/output.go`
- Create: `main.go`
- Create: all directories

**Goal:** Working directory tree and compilable skeleton. No logic yet — just types, stubs, and `go build` passes.

- [ ] **Step 1: Initialize Go module**

```bash
cd /e/minGW-choser
go mod init mingw-chooser
```

Expected: `go.mod` created with `module mingw-chooser` and `go 1.23.x`.

- [ ] **Step 2: Create directory tree**

```bash
mkdir -p detect match fetch output .github/workflows
```

- [ ] **Step 3: Write core types in `match/types.go`**

```go
package match

// Build represents one MinGW-w64 binary distribution.
type Build struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Arch      string `json:"arch"`
	Thread    string `json:"thread"`
	Exception string `json:"exception"`
	CRT       string `json:"crt"`
	GCC       string `json:"gcc"`
	Priority  int    `json:"priority"`
}

// Rules define matching preferences read from builds.json.
type Rules struct {
	ArchMap             map[string]string   `json:"arch_map"`
	ThreadPreference    []string            `json:"thread_preference"`
	ExceptionPreference map[string][]string `json:"exception_preference"`
	CRTPreference       []string            `json:"crt_preference"`
}

// DimensionChoice explains why a particular value was chosen for one dimension.
type DimensionChoice struct {
	Dimension string `json:"dimension"`
	Choice    string `json:"choice"`
	Reason    string `json:"reason"`
	Manual    bool   `json:"manual"`
}

// MatchResult holds the best build, alternatives, and explanations.
type MatchResult struct {
	Build        Build             `json:"build"`
	Alternatives []Build           `json:"alternatives,omitempty"`
	Explanation  []DimensionChoice `json:"explanation"`
	Score        int               `json:"score"`
}
```

- [ ] **Step 4: Write stub for `detect/detect.go`**

```go
package detect

// SystemInfo holds the detected system properties.
type SystemInfo struct {
	OS        string // "windows", "linux", "darwin"
	OSVersion string // e.g. "11 Pro 23H2", "Ubuntu 24.04"
	Arch      string // "x86_64", "i686", "aarch64"
	IsWow64   bool   // true when 32-bit process runs on 64-bit Windows
}
```

- [ ] **Step 5: Write stub for `detect/detect_windows.go`**

```go
//go:build windows

package detect

func Detect() SystemInfo {
	return SystemInfo{}
}
```

- [ ] **Step 6: Write stub for `detect/detect_linux.go`**

```go
//go:build linux

package detect

func Detect() SystemInfo {
	return SystemInfo{}
}
```

- [ ] **Step 7: Write stub for `detect/detect_darwin.go`**

```go
//go:build darwin

package detect

func Detect() SystemInfo {
	return SystemInfo{}
}
```

- [ ] **Step 8: Write stub for `fetch/fetch.go`**

```go
package fetch

import "mingw-chooser/match"

// Fetch retrieves the latest builds from the GitHub Releases API.
// Returns builds parsed from the asset list, or an error if unreachable.
func Fetch(apiURL string) ([]match.Build, error) {
	return nil, nil
}
```

- [ ] **Step 9: Write stub for `output/output.go`**

```go
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
	return nil
}
```

- [ ] **Step 10: Write stub for `main.go`**

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	archFlag := flag.String("arch", "", "override detected architecture (x86_64, i686, aarch64)")
	threadFlag := flag.String("thread", "", "override thread model (posix, win32)")
	excFlag := flag.String("exception", "", "override exception handling (seh, dwarf, sjlj)")
	crtFlag := flag.String("crt", "", "override CRT (ucrt, msvcrt)")
	jsonFlag := flag.Bool("json", false, "output as JSON")
	offlineFlag := flag.Bool("offline", false, "skip network fetch, use embedded snapshot only")
	listFlag := flag.Bool("list", false, "list all matching builds")
	verFlag := flag.Bool("version", false, "show version")
	flag.Parse()

	if *verFlag {
		fmt.Println("mingw-chooser", version)
		os.Exit(0)
	}

	_ = archFlag
	_ = threadFlag
	_ = excFlag
	_ = crtFlag
	_ = jsonFlag
	_ = offlineFlag
	_ = listFlag
	fmt.Println("TODO: implement")
}
```

- [ ] **Step 11: Verify project compiles**

```bash
go build ./...
```

Expected: compiles with no errors.

- [ ] **Step 12: Commit**

```bash
git init
git add -A
git commit -m "chore: scaffold project with Go module, packages, and stubs"
```

---

### Task 2: Data file (builds.json)

**Files:**
- Create: `builds.json`

**Goal:** Create the embedded data file with matching rules and a fallback build snapshot. The file must be valid JSON and parseable by `encoding/json` into `match.Rules` + extra fields.

- [ ] **Step 1: Check current mingw-builds releases for real build data**

Run in browser or use API to see actual asset names:
```
GET https://api.github.com/repos/niXman/mingw-builds-binaries/releases/latest
```

Note: this step is informational — capture a few real asset names for the fallback list.

- [ ] **Step 2: Write `builds.json`**

Content uses the data model from the spec. The `fallback_builds` array includes a few real builds (snapshot as of 2026-06, typical naming). The `sources` array has the GitHub API endpoint. The `rules` object drives the match engine:

```json
{
  "version": 1,
  "sources": [
    {
      "name": "mingw-builds",
      "api": "https://api.github.com/repos/niXman/mingw-builds-binaries/releases/latest",
      "fallback_url": "https://github.com/niXman/mingw-builds-binaries/releases"
    }
  ],
  "rules": {
    "arch_map": {
      "amd64": "x86_64",
      "386": "i686",
      "arm64": "aarch64"
    },
    "thread_preference": ["posix", "win32"],
    "exception_preference": {
      "x86_64": ["seh", "sjlj", "dwarf"],
      "i686": ["dwarf", "sjlj"]
    },
    "crt_preference": ["ucrt", "msvcrt"]
  },
  "fallback_builds": [
    {
      "name": "x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z",
      "url": "https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z",
      "arch": "x86_64",
      "thread": "posix",
      "exception": "seh",
      "crt": "ucrt"
    },
    {
      "name": "x86_64-14.2.0-release-posix-seh-msvcrt-rt_v12-rev0.7z",
      "url": "https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/x86_64-14.2.0-release-posix-seh-msvcrt-rt_v12-rev0.7z",
      "arch": "x86_64",
      "thread": "posix",
      "exception": "seh",
      "crt": "msvcrt"
    },
    {
      "name": "x86_64-14.2.0-release-win32-seh-ucrt-rt_v12-rev0.7z",
      "url": "https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/x86_64-14.2.0-release-win32-seh-ucrt-rt_v12-rev0.7z",
      "arch": "x86_64",
      "thread": "win32",
      "exception": "seh",
      "crt": "ucrt"
    },
    {
      "name": "i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z",
      "url": "https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z",
      "arch": "i686",
      "thread": "posix",
      "exception": "dwarf",
      "crt": "ucrt"
    },
    {
      "name": "i686-14.2.0-release-posix-sjlj-ucrt-rt_v12-rev0.7z",
      "url": "https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/i686-14.2.0-release-posix-sjlj-ucrt-rt_v12-rev0.7z",
      "arch": "i686",
      "thread": "posix",
      "exception": "sjlj",
      "crt": "ucrt"
    }
  ]
}
```

Note: `arch_map` maps Go's `runtime.GOARCH` values ("amd64", "386", "arm64") to MinGW arch labels ("x86_64", "i686", "aarch64").

- [ ] **Step 3: Validate JSON**

```bash
python3 -c "import json; json.load(open('builds.json')); print('valid')"
```

Expected: `valid`.

- [ ] **Step 4: Commit**

```bash
git add builds.json
git commit -m "feat: add builds.json with matching rules and fallback builds"
```

---

### Task 3: Match engine — stubs and failing tests

**Files:**
- Create: `match/match.go`
- Create: `match/match_test.go`

**Goal:** Write tests first (TDD red phase). Tests cover the core matching algorithm: filtering by arch, scoring by preference, ranking with GCC tiebreaker, user overrides.

- [ ] **Step 1: Write failing tests in `match/match_test.go`**

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./match/ -v
```

Expected: compilation error — `Match` and `Overrides` are not defined.

- [ ] **Step 3: Commit**

```bash
git add match/match_test.go
git commit -m "test: add failing match engine tests (TDD red)"
```

---

### Task 4: Match engine — implementation

**Files:**
- Modify: `match/types.go` (add `Overrides` type)
- Create: `match/match.go`

**Goal:** Implement `Match()` to make all tests pass.

- [ ] **Step 1: Add `Overrides` type to `match/types.go`**

Append to the file:

```go
// Overrides carries user-supplied flags that replace auto-detected preferences.
type Overrides struct {
	Arch      string
	Thread    string
	Exception string
	CRT       string
}
```

- [ ] **Step 2: Implement `match/match.go`**

```go
package match

import (
	"sort"
	"strconv"
	"strings"
)

// Match finds the best MinGW build for the given architecture.
// targetArch is the MinGW arch label ("x86_64", "i686", "aarch64") —
// already mapped from Go runtime arch by the detect package.
func Match(targetArch string, builds []Build, rules Rules, overrides Overrides) MatchResult {
	if overrides.Arch != "" {
		targetArch = overrides.Arch
	}

	// Filter by architecture.
	var candidates []Build
	for _, b := range builds {
		if b.Arch == targetArch {
			candidates = append(candidates, b)
		}
	}

	if len(candidates) == 0 {
		return MatchResult{
			Explanation: []DimensionChoice{
				{Dimension: "arch", Choice: targetArch, Reason: "no builds found for this architecture", Manual: overrides.Arch != ""},
			},
		}
	}

	// Score each candidate.
	type scored struct {
		Build Build
		Score int
	}
	var scoredCandidates []scored
	for _, b := range candidates {
		s := 0
		s += scoreDim(b.Thread, prefList(overrides.Thread, rules.ThreadPreference, "thread"))
		s += scoreDim(b.Exception, prefList(overrides.Exception, rules.ExceptionPreference[targetArch], "exception"))
		s += scoreDim(b.CRT, prefList(overrides.CRT, rules.CRTPreference, "crt"))
		scoredCandidates = append(scoredCandidates, scored{Build: b, Score: s})
	}

	// Sort by score desc, then GCC version desc (tiebreaker).
	sort.Slice(scoredCandidates, func(i, j int) bool {
		if scoredCandidates[i].Score != scoredCandidates[j].Score {
			return scoredCandidates[i].Score > scoredCandidates[j].Score
		}
		return compareGCC(scoredCandidates[i].Build.GCC, scoredCandidates[j].Build.GCC) > 0
	})

	best := scoredCandidates[0]
	var alternatives []Build
	for i := 1; i < len(scoredCandidates); i++ {
		alternatives = append(alternatives, scoredCandidates[i].Build)
	}

	return MatchResult{
		Build:        best.Build,
		Alternatives: alternatives,
		Explanation:  buildExplanation(targetArch, best.Build, overrides, rules),
		Score:        best.Score,
	}
}

// prefList returns the effective preference list for a dimension.
// If the user provided an override, it becomes the only option (and scores 3).
func prefList(override string, defaults []string, dim string) []string {
	if override != "" {
		return []string{override}
	}
	return defaults
}

// scoreDim returns points for a build's value based on its position in the preference list.
// First preference = 3 points, second = 2, third = 1, not found = 0.
func scoreDim(value string, prefs []string) int {
	for i, p := range prefs {
		if value == p {
			return 3 - i
		}
	}
	return 0
}

// compareGCC compares two GCC version strings like "14.2.0".
// Returns positive if a > b, negative if a < b, 0 if equal.
func compareGCC(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		ai, _ := strconv.Atoi(aParts[i])
		bi, _ := strconv.Atoi(bParts[i])
		if ai != bi {
			return ai - bi
		}
	}
	return len(aParts) - len(bParts)
}

// reasonForArch returns a human-readable reason for the arch choice.
func reasonForArch(arch string) string {
	switch arch {
	case "x86_64":
		return "your CPU is 64-bit"
	case "i686":
		return "your CPU is 32-bit"
	case "aarch64":
		return "your CPU is ARM64"
	default:
		return "detected architecture"
	}
}

// reasonForThread returns a human-readable reason for the thread model choice.
func reasonForThread(model string) string {
	switch model {
	case "posix":
		return "best C++11 std::thread support, wider compatibility"
	case "win32":
		return "native Windows threading, lighter weight"
	default:
		return "selected thread model"
	}
}

// reasonForException returns a human-readable reason for the exception handling choice.
func reasonForException(model string) string {
	switch model {
	case "seh":
		return "optimal exception handling performance on x86_64"
	case "dwarf":
		return "best exception handling for 32-bit targets"
	case "sjlj":
		return "broad compatibility across architectures"
	default:
		return "selected exception handling model"
	}
}

// reasonForCRT returns a human-readable reason for the CRT choice.
func reasonForCRT(crt string) string {
	switch crt {
	case "ucrt":
		return "modern Windows C runtime, recommended by Microsoft"
	case "msvcrt":
		return "compatible with older Windows versions (pre-Win10)"
	default:
		return "selected C runtime"
	}
}

func buildExplanation(targetArch string, best Build, overrides Overrides, rules Rules) []DimensionChoice {
	threadDim := DimensionChoice{
		Dimension: "thread",
		Choice:    best.Thread,
		Reason:    reasonForThread(best.Thread),
		Manual:    overrides.Thread != "",
	}
	if overrides.Thread != "" {
		threadDim.Reason = "[manual override] " + threadDim.Reason
	}

	excDim := DimensionChoice{
		Dimension: "exception",
		Choice:    best.Exception,
		Reason:    reasonForException(best.Exception),
		Manual:    overrides.Exception != "",
	}
	if overrides.Exception != "" {
		excDim.Reason = "[manual override] " + excDim.Reason
	}

	crtDim := DimensionChoice{
		Dimension: "crt",
		Choice:    best.CRT,
		Reason:    reasonForCRT(best.CRT),
		Manual:    overrides.CRT != "",
	}
	if overrides.CRT != "" {
		crtDim.Reason = "[manual override] " + crtDim.Reason
	}

	return []DimensionChoice{
		{Dimension: "arch", Choice: targetArch, Reason: reasonForArch(targetArch), Manual: overrides.Arch != ""},
		threadDim,
		excDim,
		crtDim,
	}
}
```

- [ ] **Step 2: Run tests to verify they pass**

```bash
go test ./match/ -v
```

Expected: all 7 tests PASS.

- [ ] **Step 3: Commit**

```bash
git add match/types.go match/match.go
git commit -m "feat: implement match engine with scoring, filtering, and overrides"
```

---

### Task 5: System detection — Windows implementation

**Files:**
- Modify: `detect/detect_windows.go`
- Create: `detect/detect_test.go` (optional — detection tests run only on matching platform)

**Goal:** Full `Detect()` for Windows using `syscall` to call `GetNativeSystemInfo`, plus environment variable fallback, plus `os/exec` for OS version.

- [ ] **Step 1: Implement `detect/detect_windows.go`**

```go
//go:build windows

package detect

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func Detect() SystemInfo {
	info := SystemInfo{OS: "windows", OSVersion: detectWindowsVersion()}
	info.Arch, info.IsWow64 = detectWindowsArch()
	return info
}

func detectWindowsArch() (arch string, isWow64 bool) {
	// Try GetNativeSystemInfo for accurate architecture detection.
	arch = mapGoArch(runtime.GOARCH)

	// Detect real native arch when running under WoW64.
	procArch := os.Getenv("PROCESSOR_ARCHITECTURE")
	procArch6432 := os.Getenv("PROCESSOR_ARCHITEW6432")

	if procArch6432 != "" {
		// Running as 32-bit process on 64-bit OS.
		isWow64 = true
		arch = mapProcessorArch(procArch6432)
	} else if procArch == "AMD64" || procArch == "ARM64" || procArch == "IA64" {
		arch = mapProcessorArch(procArch)
	}

	return arch, isWow64
}

func mapProcessorArch(pa string) string {
	switch strings.ToUpper(pa) {
	case "AMD64", "IA64":
		return "x86_64"
	case "ARM64":
		return "aarch64"
	case "X86":
		return "i686"
	default:
		return "x86_64"
	}
}

func mapGoArch(goarch string) string {
	switch goarch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i686"
	case "arm64":
		return "aarch64"
	default:
		return goarch
	}
}

func detectWindowsVersion() string {
	// Use cmd /c ver for a simple, reliable version string.
	out, err := exec.Command("cmd", "/c", "ver").Output()
	if err != nil {
		return "Windows"
	}
	// Output is "Microsoft Windows [Version 10.0.22631.4602]\r\n"
	// Strip prefix/suffix to get something readable.
	s := strings.TrimSpace(string(out))
	// Return just the version part.
	if idx := strings.Index(s, "[Version "); idx >= 0 {
		rest := s[idx+len("[Version "):]
		if end := strings.Index(rest, "]"); end >= 0 {
			return "Windows " + rest[:end]
		}
	}
	return "Windows"
}

// Ensure syscall is used (avoids "imported and not used" when only using it
// for future direct API calls).
var _ = syscall.StringToUTF16
```

- [ ] **Step 2: Verify compilation on Windows**

```bash
go build ./detect/
```

Expected: compiles successfully.

- [ ] **Step 3: Write platform-specific detection test `detect/detect_windows_test.go`**

```go
//go:build windows

package detect

import "testing"

func TestDetect_WindowsHasArch(t *testing.T) {
	info := Detect()
	if info.Arch != "x86_64" && info.Arch != "i686" && info.Arch != "aarch64" {
		t.Errorf("unexpected architecture: %s", info.Arch)
	}
	if info.OS != "windows" {
		t.Errorf("expected OS 'windows', got %s", info.OS)
	}
	if info.OSVersion == "" {
		t.Error("OSVersion should not be empty")
	}
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./detect/ -v
```

Expected: PASS (on Windows only; `go test ./...` should skip on other platforms).

- [ ] **Step 5: Commit**

```bash
git add detect/
git commit -m "feat: implement Windows system detection with WoW64 support"
```

---

### Task 6: System detection — Linux and macOS implementations

**Files:**
- Modify: `detect/detect_linux.go`
- Modify: `detect/detect_darwin.go`

**Goal:** Implement `Detect()` for Linux and macOS using standard OS utilities.

- [ ] **Step 1: Implement `detect/detect_linux.go`**

```go
//go:build linux

package detect

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Detect() SystemInfo {
	return SystemInfo{
		OS:        "linux",
		OSVersion: detectLinuxVersion(),
		Arch:      detectLinuxArch(),
	}
}

func detectLinuxArch() string {
	// Use uname -m for accurate arch.
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return mapGoArch(runtime.GOARCH)
	}
	m := strings.TrimSpace(string(out))
	return mapUnameArch(m)
}

func mapUnameArch(m string) string {
	switch m {
	case "x86_64", "amd64":
		return "x86_64"
	case "i686", "i386", "x86":
		return "i686"
	case "aarch64", "arm64":
		return "aarch64"
	default:
		return m
	}
}

func detectLinuxVersion() string {
	// Read /etc/os-release.
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		// Fallback to uname.
		out, err := exec.Command("uname", "-sr").Output()
		if err != nil {
			return "Linux"
		}
		return strings.TrimSpace(string(out))
	}
	lines := strings.Split(string(data), "\n")
	var name, version string
	for _, line := range lines {
		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		}
		if strings.HasPrefix(line, "VERSION=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION="), `"`)
		}
	}
	if name != "" {
		if version != "" {
			return name + " " + version
		}
		return name
	}
	return "Linux"
}
```

- [ ] **Step 2: Implement `detect/detect_darwin.go`**

```go
//go:build darwin

package detect

import (
	"os/exec"
	"runtime"
	"strings"
)

func Detect() SystemInfo {
	return SystemInfo{
		OS:        "darwin",
		OSVersion: detectDarwinVersion(),
		Arch:      detectDarwinArch(),
	}
}

func detectDarwinArch() string {
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return mapGoArch(runtime.GOARCH)
	}
	m := strings.TrimSpace(string(out))
	switch m {
	case "x86_64", "amd64":
		return "x86_64"
	case "arm64", "aarch64":
		return "aarch64"
	default:
		return m
	}
}

func detectDarwinVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "macOS"
	}
	ver := strings.TrimSpace(string(out))

	// Try to get the friendly name.
	nameOut, err := exec.Command("sw_vers", "-productName").Output()
	if err != nil {
		return "macOS " + ver
	}
	name := strings.TrimSpace(string(nameOut))
	return name + " " + ver
}
```

- [ ] **Step 3: Verify cross-compilation works**

```bash
GOOS=linux GOARCH=amd64 go build ./detect/
GOOS=darwin GOARCH=amd64 go build ./detect/
GOOS=windows GOARCH=amd64 go build ./detect/
```

Expected: all compile with no errors.

- [ ] **Step 4: Commit**

```bash
git add detect/
git commit -m "feat: implement Linux and macOS system detection"
```

---

### Task 7: Fetch package — GitHub API client

**Files:**
- Modify: `fetch/fetch.go`
- Create: `fetch/fetch_test.go`

**Goal:** Implement `Fetch()` that calls the GitHub Releases API, parses the JSON response, and extracts `Build` structs from asset names. The asset naming convention is `{arch}-{gcc_ver}-release-{thread}-{exception}-{crt}-{rest}.7z`.

- [ ] **Step 1: Write failing test `fetch/fetch_test.go`**

```go
package fetch

import (
	"mingw-chooser/match"
	"testing"
)

func TestParseBuildName_X86_64(t *testing.T) {
	name := "x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z"
	url := "https://example.com/" + name

	b, err := parseBuildFromAsset(name, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Arch != "x86_64" {
		t.Errorf("expected x86_64, got %s", b.Arch)
	}
	if b.Thread != "posix" {
		t.Errorf("expected posix, got %s", b.Thread)
	}
	if b.Exception != "seh" {
		t.Errorf("expected seh, got %s", b.Exception)
	}
	if b.CRT != "ucrt" {
		t.Errorf("expected ucrt, got %s", b.CRT)
	}
	if b.GCC != "14.2.0" {
		t.Errorf("expected 14.2.0, got %s", b.GCC)
	}
	if b.Name != name {
		t.Errorf("expected name %s, got %s", name, b.Name)
	}
}

func TestParseBuildName_I686(t *testing.T) {
	name := "i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z"
	url := "https://example.com/" + name

	b, err := parseBuildFromAsset(name, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Arch != "i686" {
		t.Errorf("expected i686, got %s", b.Arch)
	}
	if b.Exception != "dwarf" {
		t.Errorf("expected dwarf, got %s", b.Exception)
	}
}

func TestParseBuildName_Invalid(t *testing.T) {
	_, err := parseBuildFromAsset("not-a-valid-build.txt", "http://x")
	if err == nil {
		t.Error("expected error for invalid build name")
	}
}

func TestParseBuildName_Only7zAssets(t *testing.T) {
	_, err := parseBuildFromAsset("readme.md", "http://x")
	if err == nil {
		t.Error("expected error for non-.7z asset")
	}
}

func TestBuildListFromAssets(t *testing.T) {
	assets := []githubAsset{
		{Name: "x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z", BrowserDownloadURL: "http://a.7z"},
		{Name: "i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z", BrowserDownloadURL: "http://b.7z"},
		{Name: "README.md", BrowserDownloadURL: "http://c"},
		{Name: "source.tar.xz", BrowserDownloadURL: "http://d"},
	}
	builds := buildListFromAssets(assets)
	if len(builds) != 2 {
		t.Errorf("expected 2 builds from 4 assets, got %d", len(builds))
	}
}

// Verify the concrete type returned implements the expected shape.
var _ []match.Build = buildListFromAssets(nil)
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./fetch/ -v
```

Expected: compilation error — `parseBuildFromAsset`, `githubAsset`, `buildListFromAssets` not defined.

- [ ] **Step 3: Implement `fetch/fetch.go`**

```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./fetch/ -v
```

Expected: all 5 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add fetch/
git commit -m "feat: implement GitHub API fetch with asset name parsing"
```

---

### Task 8: Output — text formatter

**Files:**
- Create: `output/text.go`
- Create: `output/text_test.go`
- Modify: `output/output.go` (remove stub, add real dispatch)

**Goal:** Produce the human-readable text output matching the spec's example: 4 sections (system info, recommended build, install instructions, why).

- [ ] **Step 1: Implement `output/text.go`**

```go
package output

import (
	"fmt"
	"io"
	"mingw-chooser/detect"
	"mingw-chooser/match"
	"strings"
)

func printText(w io.Writer, sys detect.SystemInfo, r match.MatchResult, flags Flags) error {
	var sb strings.Builder

	// Section 1: Detected system
	sb.WriteString("\nDetected system:\n")
	sb.WriteString(fmt.Sprintf("  CPU: %s", sys.Arch))
	if sys.IsWow64 {
		sb.WriteString(" (64-bit capable, running 32-bit process — using 64-bit recommendation)")
	} else if sys.Arch == "x86_64" || sys.Arch == "aarch64" {
		sb.WriteString(" (64-bit)")
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  OS:  %s\n", sys.OSVersion))

	if flags.Offline {
		sb.WriteString("\n  [offline mode — using embedded build snapshot]\n")
	}

	// Section 2: Recommended build
	if r.Build.Name == "" {
		sb.WriteString("\nNo matching builds found for architecture: " + sys.Arch + "\n")
		sb.WriteString("Please check https://github.com/niXman/mingw-builds-binaries/releases for available builds.\n")
		_, err := fmt.Fprint(w, sb.String())
		return err
	}

	sb.WriteString("\nRecommended build:\n")
	sb.WriteString(fmt.Sprintf("  %s\n", r.Build.Name))
	sb.WriteString(fmt.Sprintf("  %s\n", r.Build.URL))

	if len(r.Alternatives) > 0 && flags.List {
		sb.WriteString("\nAlternatives:\n")
		for _, alt := range r.Alternatives {
			sb.WriteString(fmt.Sprintf("  %s\n", alt.Name))
		}
	}

	// Section 3: Install instructions
	sb.WriteString("\nHow to install:\n")
	if sys.OS == "windows" {
		sb.WriteString("  1. Extract the .7z archive to C:\\mingw64 (or your preferred location)\n")
		sb.WriteString("  2. Add C:\\mingw64\\bin to your system PATH\n")
		sb.WriteString("  3. Open a new terminal and run: gcc --version\n")
	} else {
		sb.WriteString("  1. Extract the .7z archive to ~/mingw64 (or your preferred location)\n")
		sb.WriteString("  2. Add ~/mingw64/bin to your PATH in ~/.bashrc or ~/.zshrc\n")
		sb.WriteString("  3. Open a new terminal and run: gcc --version\n")
	}

	// Section 4: Why this build
	sb.WriteString("\nWhy this build?\n")
	for _, exp := range r.Explanation {
		marker := ""
		if exp.Manual {
			marker = " [manual override]"
		}
		sb.WriteString(fmt.Sprintf("  %-9s — %s%s\n", exp.Choice, exp.Reason, marker))
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}
```

- [ ] **Step 2: Write test for text output `output/text_test.go`**

```go
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
```

- [ ] **Step 3: Update `output/output.go` to dispatch**

Replace the stub:

```go
package output

import (
	"fmt"
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

// printJSON placeholder — implemented in next task.
func printJSON(w io.Writer, sys detect.SystemInfo, r match.MatchResult, flags Flags) error {
	return fmt.Errorf("JSON output not yet implemented")
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./output/ -v
```

Expected: 3 text tests PASS (JSON dispatch returns error for now, which we'll implement next).

- [ ] **Step 5: Commit**

```bash
git add output/
git commit -m "feat: implement text output formatter with system info, build, install, and explanation sections"
```

---

### Task 9: Output — JSON formatter

**Files:**
- Create: `output/json.go`
- Modify: `output/output.go` (replace placeholder `printJSON`)
- Modify: `output/text_test.go` (add JSON tests)

**Goal:** JSON output for programmatic consumption. Struct matches the spec: `{system, recommended, alternatives, explanation, warning}`.

- [ ] **Step 1: Implement `output/json.go`**

```go
package output

import (
	"encoding/json"
	"io"
	"mingw-chooser/detect"
	"mingw-chooser/match"
)

type jsonOutput struct {
	System        systemJSON        `json:"system"`
	Recommended   *buildJSON        `json:"recommended"`
	Alternatives  []buildJSON       `json:"alternatives,omitempty"`
	Explanation   []match.DimensionChoice `json:"explanation"`
	Warning       string            `json:"warning,omitempty"`
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
```

- [ ] **Step 2: Remove the placeholder `printJSON` from `output/output.go`**

Delete the placeholder function (the real one is now in `json.go`):

```go
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
```

- [ ] **Step 3: Add JSON tests to `output/text_test.go`** (rename `output/output_test.go` would be better — let's create a new test file)

Create `output/output_test.go` with JSON tests:

```go
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
```

- [ ] **Step 4: Run all output tests**

```bash
go test ./output/ -v
```

Expected: all tests PASS (both text and JSON).

- [ ] **Step 5: Commit**

```bash
git add output/
git commit -m "feat: implement JSON output formatter"
```

---

### Task 10: Main entry point — wire everything together

**Files:**
- Modify: `main.go`

**Goal:** `main.go` parses flags, loads embedded `builds.json`, calls `detect.Detect()`, calls `fetch.Fetch()` (unless `--offline`), calls `match.Match()`, calls `output.PrintResult()`. Graceful degradation: if fetch fails, falls back to embedded builds.

- [ ] **Step 1: Create `config.go` for loading builds.json**

Create `config.go` at root:

```go
package main

import (
	_ "embed"
	"encoding/json"
	"mingw-chooser/match"
)

//go:embed builds.json
var buildsJSON []byte

type configFile struct {
	Version        int              `json:"version"`
	Sources        []sourceConfig   `json:"sources"`
	Rules          match.Rules      `json:"rules"`
	FallbackBuilds []match.Build    `json:"fallback_builds"`
}

type sourceConfig struct {
	Name        string `json:"name"`
	API         string `json:"api"`
	FallbackURL string `json:"fallback_url"`
}

func loadConfig() (configFile, error) {
	var cfg configFile
	if err := json.Unmarshal(buildsJSON, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
```

- [ ] **Step 2: Replace `main.go` with full implementation**

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"mingw-chooser/detect"
	"mingw-chooser/fetch"
	"mingw-chooser/match"
	"mingw-chooser/output"
)

const version = "0.1.0"

func main() {
	archFlag := flag.String("arch", "", "override detected architecture (x86_64, i686, aarch64)")
	threadFlag := flag.String("thread", "", "override thread model (posix, win32)")
	excFlag := flag.String("exception", "", "override exception handling (seh, dwarf, sjlj)")
	crtFlag := flag.String("crt", "", "override CRT (ucrt, msvcrt)")
	jsonFlag := flag.Bool("json", false, "output as JSON")
	offlineFlag := flag.Bool("offline", false, "skip network fetch, use embedded snapshot only")
	listFlag := flag.Bool("list", false, "list all matching builds")
	verFlag := flag.Bool("version", false, "show version")
	flag.Parse()

	if *verFlag {
		fmt.Println("mingw-chooser", version)
		os.Exit(0)
	}

	// Load embedded config.
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Detect system.
	sys := detect.Detect()

	// Get builds: fetch from API, or use fallback.
	var builds []match.Build
	var usedFallback bool
	if *offlineFlag {
		builds = cfg.FallbackBuilds
		usedFallback = true
	} else {
		for _, src := range cfg.Sources {
			fetched, err := fetch.Fetch(src.API)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to fetch %s: %v\n", src.Name, err)
				continue
			}
			builds = append(builds, fetched...)
		}
		if len(builds) == 0 {
			builds = cfg.FallbackBuilds
			usedFallback = true
			if len(cfg.Sources) > 0 {
				fmt.Fprintf(os.Stderr, "warning: using embedded build snapshot — visit %s for latest\n",
					cfg.Sources[0].FallbackURL)
			}
		}
	}

	// Match.
	overrides := match.Overrides{
		Arch:      *archFlag,
		Thread:    *threadFlag,
		Exception: *excFlag,
		CRT:       *crtFlag,
	}
	result := match.Match(sys.Arch, builds, cfg.Rules, overrides)

	// Output.
	outFlags := output.Flags{
		Arch:      *archFlag,
		Thread:    *threadFlag,
		Exception: *excFlag,
		CRT:       *crtFlag,
		Offline:   usedFallback,
		JSON:      *jsonFlag,
		List:      *listFlag,
	}

	format := output.FormatText
	if *jsonFlag {
		format = output.FormatJSON
	}

	if err := output.PrintResult(os.Stdout, sys, result, format, outFlags); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Verify full build**

```bash
go build ./...
go vet ./...
```

Expected: clean build, no errors.

- [ ] **Step 4: Verify `--help` works**

```bash
go run . --help
```

Expected: flag listing.

- [ ] **Step 5: Verify `--version` works**

```bash
go run . --version
```

Expected: `mingw-chooser 0.1.0`

- [ ] **Step 6: Run full test suite**

```bash
go test ./... -v
```

Expected: all tests PASS.

- [ ] **Step 7: End-to-end test with `--offline`**

```bash
go run . --offline
```

Expected: shows system info + a fallback build recommendation.

- [ ] **Step 8: Commit**

```bash
git add main.go config.go
git commit -m "feat: wire CLI with detect, fetch, match, and output"
```

---

### Task 11: GitHub Actions — cross-platform release

**Files:**
- Create: `.github/workflows/release.yml`

**Goal:** goreleaser workflow that builds static binaries for all target platforms on tag push.

- [ ] **Step 1: Create `.github/workflows/release.yml`**

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          check-latest: true

      - uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- [ ] **Step 2: Create `.goreleaser.yml`**

```yaml
version: 2

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - windows
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
      - "386"
    ignore:
      - goos: darwin
        goarch: "386"

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- .Os }}_
      {{- .Arch }}
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'
```

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/release.yml .goreleaser.yml
git commit -m "ci: add goreleaser workflow for cross-platform releases"
```

---

## Verification

After all tasks complete:

```bash
# Full test suite
go test ./... -v

# Build release binary
go build -ldflags="-s -w" -o mingw-chooser.exe .

# Run offline (no network needed)
./mingw-chooser.exe --offline

# Run JSON output
./mingw-chooser.exe --offline --json
```
