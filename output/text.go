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
