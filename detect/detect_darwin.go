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
