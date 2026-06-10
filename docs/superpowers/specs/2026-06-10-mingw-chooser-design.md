# MinGW Chooser — Design Spec

**Date:** 2026-06-10
**Status:** draft
**Language:** Go

## Overview

A cross-platform CLI tool that detects the user's system properties (CPU architecture, OS, runtime environment) and recommends the optimal MinGW-w64 standalone build to download. Solves the long-standing confusion around MinGW build variants (i686 vs x86_64, posix vs win32, seh vs dwarf vs sjlj, ucrt vs msvcrt).

The tool's core logic is deliberately separated from presentation so a GUI can be added later without rework.

## Goals

- Zero dependencies at runtime: a single static binary per platform.
- Works on Windows, Linux, and macOS.
- Dynamically fetches available builds from the official release API so new GCC versions are picked up automatically.
- Graceful offline fallback: ships with an embedded build snapshot; when stale or offline, guides the user to the release page.
- Data file (`builds.json`) is separate from matching logic so the community can tune preferences and add new sources (MSYS2, WinLibs, LLVM-MinGW) without code changes.
- Optional flags let users override auto-detected values.

## Non-Goals (V1)

- Interactive menus or guided wizards.
- Multi-source aggregation (MSYS2, WinLibs, etc.) — data model supports it, but initial data covers only official mingw-builds.
- Actually downloading or installing MinGW — the tool only recommends.
- Package manager integration (scoop, choco, winget).

## Architecture

```
main.go (CLI entry point)
  ├── detect/       — platform-specific system probing
  ├── match/        — scoring engine, consumes builds.json data
  ├── fetch/        — fetches latest builds from GitHub Releases API
  ├── output/       — text & JSON formatters
  └── builds.json   — embedded static data (matching rules + fallback builds)
```

### detect package

Exposes a single function `Detect() SystemInfo` implemented per-platform via build tags:

| File               | Platform |
|--------------------|----------|
| `detect_windows.go` | windows  |
| `detect_linux.go`   | linux    |
| `detect_darwin.go`  | darwin   |

**SystemInfo struct:**

```go
type SystemInfo struct {
    OS        string // "windows", "linux", "darwin"
    OSVersion string // "11 Pro 23H2", "Ubuntu 24.04"
    Arch      string // "x86_64", "i686", "aarch64"
    IsWow64   bool   // 64-bit capable CPU running a 32-bit process (Windows only)
}
```

Detection sources per platform:

- **Windows**: `GetNativeSystemInfo` + `IsWow64Process2`
- **Linux**: `/proc/cpuinfo` + `uname -m`
- **macOS**: `sysctl hw.machine` + `uname -m`

Edge case — WoW64 (32-bit process on 64-bit Windows): detect the real 64-bit capability and warn the user, recommending the 64-bit build anyway.

### builds.json

Static data file embedded at build time. Contains:

1. **Source URL** — the GitHub Releases API endpoint (or future mirrors).
2. **Matching rules** — preference orderings for each dimension.
3. **Fallback build list** — a snapshot used when offline or as a baseline when the API response can't be parsed.

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
      "x86_64": "x86_64",
      "i686": "i686",
      "aarch64": "aarch64"
    },
    "thread_preference": ["posix", "win32"],
    "exception_preference": {
      "x86_64": ["seh", "sjlj", "dwarf"],
      "i686": ["dwarf", "sjlj"]
    },
    "crt_preference": ["ucrt", "msvcrt"]
  },
  "fallback_builds": []
}
```

**Design decision — dynamic discovery vs static list:** builds.json does NOT hardcode version numbers or download URLs. Instead the `fetch` package queries the GitHub Releases API at runtime, parses the asset list, and extracts build attributes from the well-known naming convention (e.g. `x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z`). If the API is unreachable or the naming convention changes in an incompatible way, the tool falls back to the embedded snapshot and directs the user to the release page.

### match package

```go
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

type MatchResult struct {
    Build        Build
    Alternatives []Build
    Explanation  []DimensionChoice
    Score        int
}

type DimensionChoice struct {
    Dimension string // "arch", "thread", "exception", "crt"
    Choice    string // the selected value
    Reason    string // human-readable justification
    Manual    bool   // true if user overrode via flag
}
```

Matching algorithm:

1. **Filter** — keep only builds whose `arch` matches the detected (or `--arch` overridden) architecture.
2. **Score** — each dimension awards points by preference position: first preference = 3, second = 2, third = 1. Sum across dimensions.
3. **Rank** — higher total score wins. Ties broken by higher GCC version.
4. **Explain** — record why each dimension picked its value.

When the user passes a flag like `--thread win32`, the preference list for that dimension is replaced by the user's single value. The explanation marks it `[manual override]`.

### fetch package

Calls the GitHub Releases API, parses the JSON response, and extracts `Build` structs from asset names using the well-known mingw-builds naming pattern. Timeout: 10 seconds. On any error (network, parse, unexpected format) returns `nil, error` and the caller uses the fallback snapshot.

### output package

Two formatters:

**Text (default):** Colored terminal output in 4 sections:
1. Detected system info
2. Recommended build name + download URL
3. How to install (brief unpack & PATH instructions)
4. Why this build — one line per dimension

**JSON (`--json`):** Structured output for programmatic consumption / future GUI:

```json
{
  "system": {"os": "windows", "os_version": "11 Pro", "arch": "x86_64"},
  "recommended": {"name": "...", "url": "..."},
  "alternatives": [{"name": "...", "url": "..."}],
  "explanation": [
    {"dimension": "arch", "choice": "x86_64", "reason": "...", "manual": false}
  ],
  "warning": "..."  // optional, e.g. WoW64 warning
}
```

### CLI flags

```
mingw-chooser [flags]

  --arch <arch>       Override detected architecture (x86_64, i686, aarch64)
  --thread <model>    Override thread model preference (posix, win32)
  --exception <type>  Override exception handling preference (seh, dwarf, sjlj)
  --crt <type>        Override CRT preference (ucrt, msvcrt)
  --json              Output as JSON
  --offline           Skip network fetch, use embedded snapshot only
  --list              List all matching builds, not just the top recommendation
  -h, --help          Show help
  -v, --version       Show tool version
```

No arguments required — running `mingw-chooser` with zero flags is the primary use case.

### Example output

```
$ mingw-chooser

Detected system:
  CPU: x86_64 (64-bit)
  OS:  Windows 11 Pro 23H2

Recommended build:
  x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z
  https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/...

How to install:
  1. Extract the .7z archive to C:\mingw64 (or your preferred location)
  2. Add C:\mingw64\bin to your system PATH
  3. Open a new terminal and run: gcc --version

Why this build?
  x86_64  — your CPU is 64-bit
  posix   — best C++11 std::thread support, wider compatibility
  seh     — optimal exception handling performance on x86_64
  ucrt    — modern Windows C runtime, recommended by Microsoft
```

## Project structure

```
mingw-chooser/
├── main.go
├── go.mod
├── go.sum
├── builds.json              # embedded via go:embed
├── detect/
│   ├── detect.go            # SystemInfo type + Detect() interface
│   ├── detect_windows.go
│   ├── detect_linux.go
│   └── detect_darwin.go
├── match/
│   ├── match.go             # Match() engine
│   ├── match_test.go
│   └── types.go             # Build, MatchResult, Rules
├── fetch/
│   └── fetch.go             # GitHub API client + asset parser
├── output/
│   ├── output.go            # PrintResult() dispatcher
│   ├── text.go              # Text formatter
│   └── json.go              # JSON formatter
├── docs/
│   └── superpowers/
│       └── specs/
│           └── 2026-06-10-mingw-chooser-design.md
└── .github/
    └── workflows/
        └── release.yml      # goreleaser cross-compile
```

## Build & distribution

- `go build` produces a single static binary per platform.
- GitHub Actions with goreleaser cross-compiles for `windows/amd64`, `windows/arm64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`.
- Binary size target: under 10 MB.

## Future extensions (out of scope for V1)

- Additional sources in builds.json (MSYS2 packages, WinLibs, LLVM-MinGW).
- `mingw-chooser explain <dimension>` subcommand for educational output.
- GUI wrapper consuming `--json` output.
- Shell completion scripts.
- One-command download + extract (opt-in with a flag).
