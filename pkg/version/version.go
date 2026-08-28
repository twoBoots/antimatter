package version

import "fmt"

var (
	// Version holds the semantic version of Antimatter (populated via -ldflags)
	Version = "dev"
	// Commit holds the git commit hash (populated via -ldflags)
	Commit = "none"
	// Date holds the build timestamp (populated via -ldflags)
	Date = "unknown"
)

// Info contains build and version metadata
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Get returns the populated build info
func Get() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}

// String returns a human-readable representation of version info
func (i Info) String() string {
	return fmt.Sprintf("antimatter %s (commit: %s, built at: %s)", i.Version, i.Commit, i.Date)
}
