package fetch

import (
	"mingw-chooser/match"
	"testing"
)

// --- mingw-builds tests ---

func TestParseBuildName_X86_64(t *testing.T) {
	name := "x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z"
	url := "https://example.com/" + name

	b, err := parseBuildFromAsset(name, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Arch != "x86_64" {
		t.Errorf("expected x86_64, got %s", b.Arch)
	}
	if b.Thread != "posix" {
		t.Errorf("expected posix, got %s", b.Thread)
	}
	if b.Exception != "seh" {
		t.Errorf("expected seh, got %s", b.Exception)
	}
	if b.CRT != "ucrt" {
		t.Errorf("expected ucrt, got %s", b.CRT)
	}
	if b.GCC != "14.2.0" {
		t.Errorf("expected 14.2.0, got %s", b.GCC)
	}
	if b.Name != name {
		t.Errorf("expected name %s, got %s", name, b.Name)
	}
}

func TestParseBuildName_I686(t *testing.T) {
	name := "i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z"
	url := "https://example.com/" + name

	b, err := parseBuildFromAsset(name, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Arch != "i686" {
		t.Errorf("expected i686, got %s", b.Arch)
	}
	if b.Exception != "dwarf" {
		t.Errorf("expected dwarf, got %s", b.Exception)
	}
}

func TestParseBuildName_Invalid(t *testing.T) {
	_, err := parseBuildFromAsset("not-a-valid-build.txt", "http://x")
	if err == nil {
		t.Error("expected error for invalid build name")
	}
}

func TestParseBuildName_Only7zAssets(t *testing.T) {
	_, err := parseBuildFromAsset("readme.md", "http://x")
	if err == nil {
		t.Error("expected error for non-.7z asset")
	}
}

// --- WinLibs tests ---

func TestParseWinLibs_X86_64(t *testing.T) {
	name := "winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.7z"
	url := "https://example.com/" + name

	b, err := parseBuildFromAsset(name, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Arch != "x86_64" {
		t.Errorf("expected x86_64, got %s", b.Arch)
	}
	if b.Thread != "posix" {
		t.Errorf("expected posix, got %s", b.Thread)
	}
	if b.Exception != "seh" {
		t.Errorf("expected seh, got %s", b.Exception)
	}
	if b.CRT != "ucrt" {
		t.Errorf("expected ucrt, got %s", b.CRT)
	}
	if b.GCC != "16.1.0" {
		t.Errorf("expected 16.1.0, got %s", b.GCC)
	}
	if b.Name != name {
		t.Errorf("expected name %s, got %s", name, b.Name)
	}
}

func TestParseWinLibs_I686_MSVCRT(t *testing.T) {
	name := "winlibs-i686-posix-dwarf-gcc-16.1.0-mingw-w64msvcrt-14.0.0-r2.7z"
	url := "https://example.com/" + name

	b, err := parseBuildFromAsset(name, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Arch != "i686" {
		t.Errorf("expected i686, got %s", b.Arch)
	}
	if b.Exception != "dwarf" {
		t.Errorf("expected dwarf, got %s", b.Exception)
	}
	if b.CRT != "msvcrt" {
		t.Errorf("expected msvcrt, got %s", b.CRT)
	}
}

func TestParseWinLibs_Zip(t *testing.T) {
	name := "winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.zip"
	url := "https://example.com/" + name

	b, err := parseBuildFromAsset(name, url)
	if err != nil {
		t.Fatalf("unexpected error for .zip: %v", err)
	}
	if b.Arch != "x86_64" {
		t.Errorf("expected x86_64, got %s", b.Arch)
	}
}

func TestParseWinLibs_Sha256Skipped(t *testing.T) {
	_, err := parseBuildFromAsset("winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.7z.sha256", "http://x")
	if err == nil {
		t.Error("expected error for .sha256 file")
	}
}

// --- Mixed assets ---

func TestBuildListFromAssets_Mixed(t *testing.T) {
	assets := []githubAsset{
		{Name: "x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z", BrowserDownloadURL: "http://a.7z"},
		{Name: "winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.7z", BrowserDownloadURL: "http://b.7z"},
		{Name: "winlibs-i686-posix-dwarf-gcc-16.1.0-mingw-w64msvcrt-14.0.0-r2.7z", BrowserDownloadURL: "http://c.7z"},
		{Name: "README.md", BrowserDownloadURL: "http://d"},
		{Name: "source.tar.xz", BrowserDownloadURL: "http://e"},
	}
	builds := buildListFromAssets(assets)
	if len(builds) != 3 {
		t.Errorf("expected 3 builds from 5 assets (1 mingw + 2 winlibs), got %d", len(builds))
		for _, b := range builds {
			t.Logf("  %s", b.Name)
		}
	}
}

func TestBuildListFromAssets_All(t *testing.T) {
	assets := []githubAsset{
		{Name: "x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z", BrowserDownloadURL: "http://a.7z"},
		{Name: "i686-14.2.0-release-posix-dwarf-ucrt-rt_v12-rev0.7z", BrowserDownloadURL: "http://b.7z"},
		{Name: "README.md", BrowserDownloadURL: "http://c"},
		{Name: "source.tar.xz", BrowserDownloadURL: "http://d"},
	}
	builds := buildListFromAssets(assets)
	if len(builds) != 2 {
		t.Errorf("expected 2 builds from 4 assets, got %d", len(builds))
	}
}

// Verify the concrete type returned implements the expected shape.
var _ []match.Build = buildListFromAssets(nil)
