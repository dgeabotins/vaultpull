package version

import "fmt"

var (
	// Version is the current version of vaultpull.
	Version = "0.1.0"
	// Commit is the git commit hash, injected at build time.
	Commit = "none"
	// Date is the build date, injected at build time.
	Date = "unknown"
)

// Info returns a formatted version string.
func Info() string {
	return fmt.Sprintf("vaultpull v%s (commit: %s, built: %s)", Version, Commit, Date)
}

// Short returns just the semver string.
func Short() string {
	return Version
}
