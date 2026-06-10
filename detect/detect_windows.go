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
