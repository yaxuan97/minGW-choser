package fetch

import (
	"mingw-chooser/match"
	"testing"
)

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

func TestBuildListFromAssets(t *testing.T) {
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
