package tui

import "github.com/charmbracelet/lipgloss"

var (
	// titleStyle is the style for the application title
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	// emptyStyle is the style for empty state messages
	emptyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	// messageStyle is the style for success messages
	messageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	// errorStyle is the style for error messages
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	// helpStyle is the style for help text
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	// selectedStyle is the style for selected items
	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	// todoItemStyle is the style for regular todo items
	todoItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// completedItemStyle is the style for completed todo items
	completedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Strikethrough(true)

	// highPriorityStyle is the style for high priority indicator
	highPriorityStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	// mediumPriorityStyle is the style for medium priority indicator
	mediumPriorityStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("220"))

	// lowPriorityStyle is the style for low priority indicator
	lowPriorityStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("82"))

	// bannerStyle is the style for the startup banner
	bannerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06c775")).
		Bold(true).
		Align(lipgloss.Center)

	// bannerSubtitleStyle is the style for the banner subtitle
	bannerSubtitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Italic(true).
		Align(lipgloss.Center)

	// bannerVersionStyle is the style for the version info
	bannerVersionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	// bannerPromptStyle is the style for the "press any key" prompt
	bannerPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("246")).
		Italic(true).
		Align(lipgloss.Center).
		MarginTop(2)
)
