package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/syeeel/koto-cli-go/internal/config"
	"github.com/syeeel/koto-cli-go/internal/repository"
	"github.com/syeeel/koto-cli-go/internal/service"
	"github.com/syeeel/koto-cli-go/internal/tui"
)

// Version information (set by GoReleaser via ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version information for TUI
	tui.Version = version
	tui.CommitSHA = commit
	tui.BuildDate = date

	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit:  %s\n", commit)
		fmt.Printf("Built:   %s\n", date)
		return
	}
	// Get configuration
	cfg, err := config.GetDefaultConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: failed to get configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	repo, err := repository.NewSQLiteRepository(cfg.DBPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to close database: %v\n", err)
		}
	}()

	// Initialize service
	svc := service.NewTodoService(repo)

	// Create TUI model
	model := tui.NewModel(svc)

	// Start the application
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
