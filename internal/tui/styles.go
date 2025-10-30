package tui

import "github.com/charmbracelet/lipgloss"

// MinTerminalWidth is the minimum required terminal width
const MinTerminalWidth = 100

var (
	// Theme colors (transparent backgrounds)
	fgDefault   = lipgloss.Color("#cdd6f4") // Light text
	fgSelected  = lipgloss.Color("#39ff14") // Neon green for selected text
	fgDim       = lipgloss.Color("#6c7086") // Dimmed text
	fgCompleted = lipgloss.Color("#585b70") // Completed items
	accentGreen = lipgloss.Color("#a6e3a1") // Accent color
	accentRed   = lipgloss.Color("#f38ba8") // Error color

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
			Bold(true).
			Width(116)

	// todoItemStyle is the style for regular todo items (transparent background)
	todoItemStyle = lipgloss.NewStyle().
			Foreground(fgDefault)

	// completedItemStyle is the style for completed todo items (transparent background)
	completedItemStyle = lipgloss.NewStyle().
				Foreground(fgCompleted).
				Strikethrough(true)

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
			Foreground(lipgloss.Color("213")).
			Bold(true).
			Underline(true)
)

// DynamicWidths holds calculated widths for responsive layout
type DynamicWidths struct {
	// Main list view column widths
	NoCol         int
	TitleCol      int
	PriorityCol   int
	WorkTimeCol   int
	CreatedCol    int
	TotalListRow  int

	// Detail view widths
	DetailBox     int
	DetailColumn  int

	// Pomodoro view widths
	PomodoroBox   int
	ProgressBar   int

	// Edit view widths
	EditInput     int

	// General
	ContentWidth  int
}

// calculateDynamicWidths calculates responsive widths based on terminal width
func calculateDynamicWidths(termWidth int) DynamicWidths {
	// Reserve space for margins and borders
	const (
		marginBuffer  = 4
		borderBuffer  = 2
	)

	contentWidth := termWidth - marginBuffer

	// Main list view - proportional column widths
	// Format: No(5) | Title(flexible) | Priority(12) | WorkTime(12) | Created(13)
	fixedColsWidth := 5 + 12 + 12 + 13 + 8 // 8 for spacing between columns
	titleWidth := contentWidth - fixedColsWidth
	if titleWidth < 20 {
		titleWidth = 20 // Minimum title width
	}
	totalRowWidth := contentWidth

	// Detail view - use most of terminal width
	detailBoxWidth := contentWidth - 6
	if detailBoxWidth < 60 {
		detailBoxWidth = 60
	}
	detailColumnWidth := (contentWidth - 12) / 3
	if detailColumnWidth < 20 {
		detailColumnWidth = 20
	}

	// Pomodoro view - centered content
	pomodoroBoxWidth := contentWidth - 26
	if pomodoroBoxWidth < 50 {
		pomodoroBoxWidth = 50
	}
	progressBarWidth := contentWidth - 36
	if progressBarWidth < 40 {
		progressBarWidth = 40
	}

	// Edit view - input width
	editInputWidth := contentWidth - 16
	if editInputWidth < 40 {
		editInputWidth = 40
	}

	return DynamicWidths{
		NoCol:         5,
		TitleCol:      titleWidth,
		PriorityCol:   12,
		WorkTimeCol:   12,
		CreatedCol:    13,
		TotalListRow:  totalRowWidth,
		DetailBox:     detailBoxWidth,
		DetailColumn:  detailColumnWidth,
		PomodoroBox:   pomodoroBoxWidth,
		ProgressBar:   progressBarWidth,
		EditInput:     editInputWidth,
		ContentWidth:  contentWidth,
	}
}

// createResponsiveBoxStyle creates a box style with dynamic width
func createResponsiveBoxStyle(width int, borderStyle lipgloss.Border, borderColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(width)
}
