# Development Guide

This document describes the architecture, design decisions, and development workflow for MinGW Chooser. It's intended for anyone who wants to understand how the project works or contribute to it.

## Project Overview

MinGW Chooser is a cross-platform CLI tool (Go) that detects the user's system and recommends the best MinGW-w64 standalone build to download. Single static binary, zero runtime dependencies.

## Quick Start

```bash
# Build
go build -o mingw-chooser.exe .

# Run tests
go test ./...

# Run a single test
go test ./match/ -run TestMatch -v

# Cross-compile for all targets
GOOS=windows GOARCH=amd64 go build -o mingw-chooser-win64.exe .
GOOS=linux   GOARCH=amd64 go build -o mingw-chooser-linux64 .
GOOS=darwin  GOARCH=amd64 go build -o mingw-chooser-darwin64 .
GOOS=darwin  GOARCH=arm64 go build -o mingw-chooser-darwin-arm64 .

# Lint (requires golangci-lint)
golangci-lint run

# Run with flags
go run . --json
go run . --offline
go run . --arch i686 --thread posix
```

## Architecture

Four independent packages, orchestrated by `main.go`:

```
main.go
  ├── detect/    Platform-specific system probing
  ├── fetch/     GitHub Releases API client
  ├── match/     Scoring engine
  ├── output/    Text & JSON formatters
  └── builds.json (embedded)  Rules + fallback snapshot
```

| Package | Responsibility | Key export |
|---------|---------------|------------|
| `detect/` | Platform-specific system probing (build tags: windows/linux/darwin) | `Detect() SystemInfo` |
| `fetch/` | GitHub Releases API client, parses asset names into Build structs | Returns `[]Build` or error |
| `match/` | Scoring engine — filters builds by arch, scores by preference order, ranks | `Match(targetArch, []Build, Rules, Overrides) MatchResult` |
| `output/` | Text (default) and JSON (`--json`) formatters | `PrintResult(w, SystemInfo, MatchResult, Format, Flags)` |

`builds.json` is embedded at build time via `go:embed`. It contains matching rules, source configurations, and a fallback build snapshot. The file is separate from match logic so preferences can be tuned without code changes.

## Data Flow

1. `main.go` parses flags → loads embedded `builds.json` → calls `detect.Detect()`
2. If `--offline` is **not** set, calls `fetch.Fetch()` for each source in `builds.json`; on failure, uses embedded fallback
3. Calls `match.Match(systemInfo.Arch, builds, rules, overrides)` → filters by arch, scores each dimension by preference order
4. Calls `output.PrintResult()` with the chosen format (text or JSON)

## Matching Algorithm

1. **Filter** — keep only builds matching the target architecture
2. **Score** — each dimension awards points by preference position (1st = 3 pts, 2nd = 2 pts, 3rd = 1 pt). Sum across thread, exception, and CRT dimensions
3. **Rank** — higher score wins. Ties broken in order: GCC version → source priority
4. **Explain** — records why each dimension chose its value

### Why Source Priority

When two sources offer builds with identical dimensions and the same GCC version, source priority acts as a tiebreaker. Currently WinLibs (priority 10) is preferred over mingw-builds (priority 0) — it tends to bundle newer CRT libraries and extra tools that benefit Windows users. You can adjust these priorities in `builds.json`.

## Sources

The tool fetches from multiple GitHub Releases endpoints:

| Source | API | Priority |
|--------|-----|----------|
| [mingw-builds](https://github.com/niXman/mingw-builds-binaries) | `.../releases/latest` | 0 |
| [WinLibs](https://github.com/brechtsanders/winlibs_mingw) | `.../releases/latest` | 10 |

Each source has its own asset naming convention. The `fetch` package uses separate regex patterns to parse them:

```
mingw-builds:  x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z
WinLibs:       winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.7z
```

To add a new source, you need to:
1. Add an entry in `builds.json` → `sources` with its API URL and priority
2. Add a regex pattern and parser in `fetch/fetch.go`
3. Add fallback builds under `fallback_builds`

## Key Design Decisions

- **Dynamic discovery over static build list** — `builds.json` does not hardcode version numbers or download URLs. The tool fetches the latest release from GitHub API at runtime and parses asset names. New GCC versions are picked up automatically.
- **Graceful offline** — When the API is unreachable, falls back to an embedded build snapshot and directs users to the release page.
- **User overrides** — `--arch`, `--thread`, `--exception`, `--crt` flags replace auto-detected values. Overridden dimensions are marked `[manual]` in explanations.
- **Scoring** — Each dimension awards points by preference position. Ties broken by GCC version, then source priority.
- **Multi-source** — The data model supports multiple build sources. Scoring works uniformly across all of them.

## Edge Cases

- **WoW64** — 32-bit process on 64-bit Windows: detect the real 64-bit capability, warn, recommend 64-bit build
- **ARM64 Windows** — check if ARM64 builds exist in the fetched set; if not, suggest x86_64 cross-compilation
- **Naming convention drift** — if the GitHub API response can't be parsed, fall back to embedded snapshot and alert the user
- **All sources fail** — if every fetch fails, use the embedded fallback snapshot

## Project Structure

```
mingw-chooser/
├── main.go              CLI entry point + flag parsing
├── config.go            Embedded `builds.json` loader
├── builds.json          Matching rules, source configs, fallback builds
├── detect/
│   ├── detect.go        SystemInfo type
│   ├── detect_windows.go  Windows implementation (build tag)
│   ├── detect_linux.go    Linux implementation (build tag)
│   ├── detect_darwin.go   macOS implementation (build tag)
│   └── detect_windows_test.go
├── fetch/
│   ├── fetch.go         API client + dual-format asset parser
│   └── fetch_test.go
├── match/
│   ├── types.go         Build, Rules, Overrides, MatchResult
│   ├── match.go         Match() + scoring + sort + explanations
│   └── match_test.go
├── output/
│   ├── output.go        PrintResult() dispatcher
│   ├── text.go          Human-readable formatter
│   ├── json.go          JSON formatter
│   ├── text_test.go
│   └── json_test.go
├── .github/workflows/
│   └── release.yml      goreleaser cross-platform CI
└── .goreleaser.yml
```

## Testing

All packages have tests. Run the full suite before committing:

```bash
go test ./... -v

# Coverage
go test ./... -cover
```

Tests are package-scoped — each package's tests verify its own behavior. The match engine tests load `builds.json` at test time to stay in sync with the real preference data.
