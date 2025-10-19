package tui

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
	return "v1.0.0"
}
