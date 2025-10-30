package tui

import "fmt"

// Version information (set via ldflags during build)
var (
	Version   = "dev"      // Version number
	CommitSHA = "none"     // Git commit SHA
	BuildDate = "unknown"  // Build date
)

// GetBanner returns the KOTO CLI ASCII art banner
func GetBanner() string {
	return `
██╗  ██╗ ██████╗ ████████╗ ██████╗      ██████╗██╗     ██╗
██║ ██╔╝██╔═══██╗╚══██╔══╝██╔═══██╗    ██╔════╝██║     ██║
█████╔╝ ██║   ██║   ██║   ██║   ██║    ██║     ██║     ██║
██╔═██╗ ██║   ██║   ██║   ██║   ██║    ██║     ██║     ██║
██║  ██╗╚██████╔╝   ██║   ╚██████╔╝    ╚██████╗███████╗██║
╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝      ╚═════╝╚══════╝╚═╝
`
}

// GetSubtitle returns the subtitle for the banner
func GetSubtitle() string {
	return "✨ Your Beautiful Terminal ToDo Manager ✨"
}

// GetVersion returns the version string
func GetVersion() string {
	return fmt.Sprintf("koto version %s\n  commit: %s\n  built:  %s", Version, CommitSHA, BuildDate)
}
