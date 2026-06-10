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
