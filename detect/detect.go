package detect

// SystemInfo holds the detected system properties.
type SystemInfo struct {
	OS        string // "windows", "linux", "darwin"
	OSVersion string // e.g. "11 Pro 23H2", "Ubuntu 24.04"
	Arch      string // "x86_64", "i686", "aarch64"
	IsWow64   bool   // true when 32-bit process runs on 64-bit Windows
}
