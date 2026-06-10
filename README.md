# MinGW Chooser

[中文版](README.zh-CN.md)

> **⚠️ This project is under active development.** The core functionality works, but APIs, flags, and output formats may change before v1.0. Feedback and contributions are welcome!

A cross-platform CLI tool that detects your system and recommends the best MinGW-w64 build to download. No more confusion between i686 vs x86_64, posix vs win32, seh vs dwarf vs sjlj, ucrt vs msvcrt.

## Quick Start

```bash
# Download the latest binary from Releases, or build from source:
go install github.com/yourusername/mingw-chooser@latest

# Run with zero flags — auto-detect and recommend:
mingw-chooser

# Offline mode (no network):
mingw-chooser --offline

# JSON output for scripting:
mingw-chooser --json
```

## What It Does

Running `mingw-chooser` with no flags:

1. **Detects** your CPU architecture, OS, and whether you're running under WoW64
2. **Fetches** the latest available builds from [mingw-builds](https://github.com/niXman/mingw-builds-binaries) and [WinLibs](https://github.com/brechtsanders/winlibs_mingw)
3. **Scores** each build by how well it matches your system (posix > win32, seh > dwarf on x86_64, ucrt > msvcrt)
4. **Recommends** the single best build with a download link, install instructions, and an explanation of every choice

## Example Output

```
$ mingw-chooser

Detected system:
  CPU: x86_64 (64-bit)
  OS:  Windows 11 Pro

Recommended build:
  winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.7z
  https://github.com/brechtsanders/winlibs_mingw/releases/...

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

## Flags

| Flag | Description |
|------|-------------|
| `--arch <arch>` | Override architecture (`x86_64`, `i686`, `aarch64`) |
| `--thread <model>` | Override thread model (`posix`, `win32`) |
| `--exception <type>` | Override exception handling (`seh`, `dwarf`, `sjlj`) |
| `--crt <type>` | Override CRT (`ucrt`, `msvcrt`) |
| `--json` | Output as JSON |
| `--offline` | Use embedded build snapshot only (no network) |
| `--list` | Show all matching builds, not just the top pick |
| `--version` | Show version |

## JSON Output

```json
{
  "system": {"os": "windows", "os_version": "Windows 11 Pro", "arch": "x86_64"},
  "recommended": {"name": "winlibs-x86_64-...", "url": "https://..."},
  "alternatives": [...],
  "explanation": [
    {"dimension": "arch", "choice": "x86_64", "reason": "your CPU is 64-bit", "manual": false}
  ],
  "warning": null
}
```

## How It Works

```
main.go
  ├── detect/    Platform detection (Windows/Linux/macOS)
  ├── fetch/     GitHub Releases API client (mingw-builds + WinLibs)
  ├── match/     Scoring engine — filters, scores, ranks
  ├── output/    Text & JSON formatters
  └── builds.json (embedded)  Rules + fallback snapshot
```

### Matching Algorithm

1. **Filter** — keep builds matching the target architecture
2. **Score** — each dimension awards points by preference position (1st = 3 pts, 2nd = 2 pts, 3rd = 1 pt)
3. **Rank** — higher score wins. Ties broken by: GCC version → source priority
4. **Explain** — every dimension choice is explained

### Sources

| Source | Priority | Notes |
|--------|----------|-------|
| [WinLibs](https://winlibs.com/) | High | Frequently updated, bundles extra libraries |
| [mingw-builds](https://github.com/niXman/mingw-builds-binaries) | Base | Official standalone builds |

The tool fetches from **both** sources, scores all builds together, and picks the best one. WinLibs gets a slight priority advantage when specs are otherwise identical — it tends to be more up-to-date for Windows users.

### Edge Cases Handled

- **WoW64** — 32-bit process on 64-bit Windows? Detects the real capability, warns, recommends 64-bit
- **ARM64 Windows** — recommends x86_64 cross-compilation if no native ARM64 build exists
- **Offline** — falls back to embedded build snapshot, directs user to release page
- **Naming drift** — if API response can't be parsed, falls back gracefully

## Build from Source

```bash
# Requires Go 1.23+
git clone https://github.com/yourusername/mingw-chooser.git
cd mingw-chooser
go build -o mingw-chooser .

# Cross-compile
GOOS=windows GOARCH=amd64 go build -o mingw-chooser.exe .
GOOS=linux   GOARCH=amd64 go build -o mingw-chooser .
GOOS=darwin  GOARCH=amd64 go build -o mingw-chooser .
```

Zero external dependencies — standard library only.

## Project Structure

```
mingw-chooser/
├── main.go              CLI entry point
├── config.go            Embedded config loader
├── builds.json          Matching rules + fallback builds
├── detect/              Platform detection (build tags)
├── fetch/               GitHub API client + asset parser
├── match/               Scoring engine + tests
└── output/              Text & JSON formatters
```

## Contributing

Contributions are welcome! Here's how to get started:

1. **Fork** the repository and clone it locally
2. **Create a branch** for your change (`git checkout -b feat/my-feature`)
3. **Make your changes** — follow existing code style and add tests
4. **Run the test suite** — `go test ./...` must pass
5. **Commit** with a [conventional commit](https://www.conventionalcommits.org/) message
6. **Open a pull request** with a clear description of what you changed and why

### What to Work On

- **New sources** — add MSYS2, LLVM-MinGW, or other MinGW distributions to `builds.json` and `fetch/`
- **New platforms** — improve detection for niche environments (Cygwin, MSYS2, WSL)
- **Better recommendations** — tune scoring rules in `builds.json` based on real-world experience
- **Shell completions** — write bash/zsh/fish completion scripts
- **One-click install** — add a `--download` flag to download and extract the recommended build
- **GUI wrapper** — build a desktop frontend that consumes `--json` output

### Development Setup

```bash
git clone https://github.com/yourusername/mingw-chooser.git
cd mingw-chooser
go test ./...        # verify everything passes
go build .           # build the binary
```

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed architecture, design decisions, and development workflow.

## License

MIT
