package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Theme colors (transparent backgrounds)
	fgDefault       = lipgloss.Color("#cdd6f4")   // Light text
	fgHeader        = lipgloss.Color("#f5e0dc")   // Header text (lighter)
	fgSelected      = lipgloss.Color("#39ff14")   // Neon green for selected text
	fgDim           = lipgloss.Color("#6c7086")   // Dimmed text
	fgCompleted     = lipgloss.Color("#585b70")   // Completed items
	accentGreen     = lipgloss.Color("#a6e3a1")   // Accent color
	accentRed       = lipgloss.Color("#f38ba8")   // Error color
	separatorColor  = lipgloss.Color("#313244")   // Subtle separator

	// titleStyle is the style for the application title (transparent background)
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(accentGreen).
		Padding(0, 1).
		MarginBottom(1)

	// emptyStyle is the style for empty state messages (transparent background)
	emptyStyle = lipgloss.NewStyle().
		Foreground(fgDim).
		Italic(true)

	// messageStyle is the style for success messages (transparent background)
	messageStyle = lipgloss.NewStyle().
		Foreground(accentGreen).
		Bold(true)

	// errorStyle is the style for error messages (transparent background)
	errorStyle = lipgloss.NewStyle().
		Foreground(accentRed).
		Bold(true)

	// helpStyle is the style for help text (transparent background)
	helpStyle = lipgloss.NewStyle().
		Foreground(fgDim).
		MarginTop(1)

	// selectedStyle is the style for selected items (neon green background, dark text)
	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")).
		Background(fgSelected).
		Bold(true)

	// todoItemStyle is the style for regular todo items (transparent background)
	todoItemStyle = lipgloss.NewStyle().
		Foreground(fgDefault)

	// todoItemAltStyle is the style for alternate row todo items (transparent background)
	todoItemAltStyle = lipgloss.NewStyle().
		Foreground(fgDefault)

	// completedItemStyle is the style for completed todo items (transparent background)
	completedItemStyle = lipgloss.NewStyle().
		Foreground(fgCompleted).
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

	// bannerTodoBoxStyle is the style for the todo list box on the banner screen
	bannerTodoBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#06c775")).
		Padding(1, 2).
		Width(40)

	// bannerTodoTitleStyle is the style for the todo box title
	bannerTodoTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06c775")).
		Bold(true).
		Align(lipgloss.Center)

	// bannerTodoItemStyle is the style for todo items in the banner
	bannerTodoItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// bannerTodoNumberStyle is the style for todo item numbers
	bannerTodoNumberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06c775")).
		Bold(true)

	// headerStyle is the style for table headers (transparent background)
	headerStyle = lipgloss.NewStyle().
		Foreground(fgHeader).
		Bold(true).
		Underline(true)

	// separatorStyle is the style for table separators (transparent background)
	separatorStyle = lipgloss.NewStyle().
		Foreground(separatorColor)

	// inputStyle is the style for input field (transparent background)
	inputStyle = lipgloss.NewStyle().
		Foreground(fgDefault)
)
