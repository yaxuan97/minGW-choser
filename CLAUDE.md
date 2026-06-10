# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

MinGW Chooser — a cross-platform CLI tool (Go) that detects the user's system and recommends the correct MinGW-w64 standalone build to download. Single static binary, zero runtime dependencies.

## Commands

```bash
# Build
go build -o mingw-chooser.exe .

# Run tests
go test ./...

# Run a single test
go test ./match/ -run TestMatch -v

# Cross-compile for all targets
GOOS=windows GOARCH=amd64 go build -o dist/mingw-chooser-win64.exe .
GOOS=linux   GOARCH=amd64 go build -o dist/mingw-chooser-linux64 .
GOOS=darwin  GOARCH=amd64 go build -o dist/mingw-chooser-darwin64 .
GOOS=darwin  GOARCH=arm64 go build -o dist/mingw-chooser-darwin-arm64 .

# Lint (requires golangci-lint)
golangci-lint run

# Run with flags
go run . --json
go run . --offline
go run . --arch i686 --thread posix
```

## Architecture

Four independent packages, orchestrated by `main.go`:

| Package | Responsibility | Key export |
|---------|---------------|------------|
| `detect/` | Platform-specific system probing (build tags: windows/linux/darwin) | `Detect() SystemInfo` |
| `fetch/` | GitHub Releases API client, parses asset names into Build structs | Returns `[]Build` or error |
| `match/` | Scoring engine — filters builds by arch, scores by preference order, ranks | `Match(SystemInfo, []Build, Rules) MatchResult` |
| `output/` | Text (default) and JSON (`--json`) formatters | `PrintResult(w, MatchResult, Format, Flags)` |

`builds.json` is embedded at build time via `go:embed`. It contains matching rules + a fallback build snapshot. The file is separate from match logic so preferences can be tuned without code changes and new sources (MSYS2, WinLibs) can be added later.

## Data flow

1. `main.go` parses flags → loads embedded `builds.json` → calls `detect.Detect()`
2. If `--offline` is NOT set, calls `fetch.Fetch()` for latest builds from GitHub API; on failure, uses embedded fallback
3. Calls `match.Match(systemInfo, builds, rules)` → filters by arch, scores each dimension by preference order
4. Calls `output.PrintResult()` with the chosen format (text or JSON)

## Key design decisions

- **Dynamic discovery over static build list**: `builds.json` does NOT hardcode version numbers or download URLs. The tool fetches the latest release from the GitHub API at runtime and parses asset names. This means new GCC versions are picked up automatically without updating the tool.
- **Graceful offline**: When the API is unreachable, falls back to an embedded build snapshot and directs users to the release page.
- **User overrides**: `--arch`, `--thread`, `--exception`, `--crt` flags replace auto-detected values. Overridden dimensions are marked `[manual]` in explanations.
- **Scoring**: Each dimension awards points by preference position (first choice = 3, second = 2, third = 1). Ties broken by higher GCC version.
- **V1 scope**: Only official mingw-builds (niXman/mingw-builds-binaries). The data model supports adding MSYS2, WinLibs, LLVM-MinGW later.

## Edge cases to handle

- **WoW64** (32-bit process on 64-bit Windows): detect the real 64-bit capability, warn, recommend 64-bit build.
- **ARM64 Windows**: check if ARM64 builds exist in the fetched set; if not, suggest x86_64 cross-compilation.
- **Naming convention drift**: if the GitHub API response can't be parsed, fall back to embedded snapshot and alert the user.
