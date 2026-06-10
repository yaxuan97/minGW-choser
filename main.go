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
