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
