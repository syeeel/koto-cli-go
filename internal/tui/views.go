package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/syeeel/koto-cli-go/internal/model"
)

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	switch m.viewMode {
	case ViewModeBanner:
		return m.renderBannerView()
	case ViewModeHelp:
		return m.renderHelpView()
	case ViewModeAddTodo:
		return m.renderAddTodoView()
	case ViewModeEditTodo:
		return m.renderEditTodoView()
	case ViewModePomodoro:
		return m.renderPomodoroView()
	default:
		return m.renderListView()
	}
}

// renderListView renders the main todo list view
func (m Model) renderListView() string {
	var s strings.Builder

	// Title with dark background
	s.WriteString(titleStyle.Render(" 📝 koto - ToDo Manager "))
	s.WriteString("\n\n")

	// Todo list
	if len(m.todos) == 0 {
		s.WriteString(emptyStyle.Render("  No todos yet. Use /add to create your first todo!  "))
		s.WriteString("\n")
	} else {
		// Header with fixed widths: No.(4), Title(35), Description(35), Total time(10), Create Date(11)
		headerNo := padStringToWidth("No.", 4)
		headerTitle := padStringToWidth("Title", 35)
		headerDesc := padStringToWidth("Description", 35)
		headerTime := padStringToWidth("Total time", 10)
		headerDate := padStringToWidth("Create Date", 11)
		header := fmt.Sprintf(" %s   %s   %s   %s   %s ", headerNo, headerTitle, headerDesc, headerTime, headerDate)
		s.WriteString(headerStyle.Render(header))
		s.WriteString("\n")

		// Todo items (no separator lines between items for cleaner look)
		for i, todo := range m.todos {
			s.WriteString(m.renderTodoItem(i, todo))
			s.WriteString("\n")
		}
	}

	// Input field
	s.WriteString("\n")
	s.WriteString(m.input.View())
	s.WriteString("\n")

	// Status messages
	if m.message != "" {
		s.WriteString("\n")
		s.WriteString(messageStyle.Render(m.message))
		s.WriteString("\n")
	}

	if m.err != nil {
		s.WriteString("\n")
		s.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		s.WriteString("\n")
	}

	// Help text
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("Commands: /add, /list, /done, /delete, /edit, /pomo, /help | Navigate: ↑/↓ or j/k | Help: ? | Quit: /exit or Ctrl+C"))

	return s.String()
}

// renderTodoItem renders a single todo item in table format
func (m Model) renderTodoItem(index int, todo *model.Todo) string {
	// No. (ID) - width: 4
	no := fmt.Sprintf("%d", todo.ID)
	no = padStringToWidth(no, 4)

	// Title - width: 35 (display width, not character count)
	title := truncateStringByWidth(todo.Title, 35)
	title = padStringToWidth(title, 35)

	// Description - width: 35 (display width, not character count)
	desc := truncateStringByWidth(todo.Description, 35)
	desc = padStringToWidth(desc, 35)

	// Total time - width: 10
	totalTime := todo.GetWorkDurationFormatted()
	if totalTime == "" {
		totalTime = "-"
	}
	totalTime = padStringToWidth(totalTime, 10)

	// Create Date (format: YYYY-MM-DD) - width: 11
	createDate := todo.CreatedAt.Format("2006-01-02")
	createDate = padStringToWidth(createDate, 11)

	// Build the row with spacing (no vertical separators for cleaner look)
	row := fmt.Sprintf(" %s   %s   %s   %s   %s ", no, title, desc, totalTime, createDate)

	// Apply cursor style if selected (neon green text, transparent background)
	if m.cursor == index {
		return selectedStyle.Render(row)
	}

	// Apply completed style if todo is completed
	if todo.IsCompleted() {
		return completedItemStyle.Render(row)
	}

	// No alternating backgrounds - transparent background for all rows
	return todoItemStyle.Render(row)
}

// truncateString truncates a string to maxLength characters (deprecated, use truncateStringByWidth)
func truncateString(s string, maxLength int) string {
	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}
	return string(runes[:maxLength-3]) + "..."
}

// truncateStringByWidth truncates a string based on display width
// considering that fullwidth characters (Japanese, Chinese, etc.) take 2 cells
func truncateStringByWidth(s string, maxWidth int) string {
	width := runewidth.StringWidth(s)
	if width <= maxWidth {
		return s
	}

	// Build string until we reach maxWidth - 3 (for "...")
	targetWidth := maxWidth - 3
	result := ""
	currentWidth := 0

	for _, r := range s {
		charWidth := runewidth.RuneWidth(r)
		if currentWidth+charWidth > targetWidth {
			break
		}
		result += string(r)
		currentWidth += charWidth
	}

	return result + "..."
}

// padStringToWidth pads a string to a specific display width
func padStringToWidth(s string, width int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

// renderPriority renders the priority indicator
func (m Model) renderPriority(priority model.Priority) string {
	switch priority {
	case model.PriorityHigh:
		return highPriorityStyle.Render("🔴")
	case model.PriorityMedium:
		return mediumPriorityStyle.Render("🟡")
	case model.PriorityLow:
		return lowPriorityStyle.Render("🟢")
	default:
		return "⚪"
	}
}

// renderHelpContent generates the help content for viewport
func (m Model) renderHelpContent() string {
	var s strings.Builder

	// Title with dark background
	title := titleStyle.Render(" 📖 koto - Help ")
	s.WriteString(title)
	s.WriteString("\n\n")

	// Scroll hint
	scrollHintStyle := lipgloss.NewStyle().Foreground(accentGreen).Italic(true)
	s.WriteString(scrollHintStyle.Render("  💡 You can scroll this page using ↑/↓ or j/k keys  "))
	s.WriteString("\n\n")

	// Commands header with style
	cmdHeader := headerStyle.Render(" COMMANDS ")
	s.WriteString(cmdHeader)
	s.WriteString("\n\n")

	commands := []struct {
		command string
		desc    string
		example string
	}{
		{"/add", "Add a new todo (interactive)", "/add"},
		{"", "  → Step 1: Enter title", ""},
		{"", "  → Step 2: Enter description (optional)", ""},
		{"", "", ""},
		{"/list", "List all todos", "/list"},
		{"/list --status=<pending|completed>", "List by status", "/list --status=pending"},
		{"", "", ""},
		{"/done <id>", "Mark todo as completed", "/done 1"},
		{"/delete <id>", "Delete a todo", "/delete 2"},
		{"/edit <id>", "Edit a todo (interactive)", "/edit 1"},
		{"", "  → Step 1: Edit title", ""},
		{"", "  → Step 2: Edit description (optional)", ""},
		{"", "", ""},
		{"/pomo [id]", "Start a 25-minute Pomodoro timer", "/pomo"},
		{"", "  → General timer (no task)", "/pomo"},
		{"", "  → Task-specific timer (records time)", "/pomo 1"},
		{"", "", ""},
		{"/export [filepath]", "Export todos to JSON", "/export ~/todos.json"},
		{"/import <filepath>", "Import todos from JSON", "/import ~/todos.json"},
		{"", "", ""},
		{"/help", "Show this help screen", "/help"},
		{"/exit", "Quit the application", "/exit"},
	}

	// Command items with transparent background
	cmdStyle := lipgloss.NewStyle().Foreground(accentGreen).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(fgDefault)
	exampleLabelStyle := lipgloss.NewStyle().Foreground(fgDim).Italic(true)

	for _, cmd := range commands {
		if cmd.command == "" {
			s.WriteString("\n")
			continue
		}
		s.WriteString(fmt.Sprintf("  %s  %s\n",
			cmdStyle.Render(fmt.Sprintf("%-45s", cmd.command)),
			descStyle.Render(cmd.desc)))
		if cmd.example != "" {
			s.WriteString(fmt.Sprintf("    %s %s\n",
				exampleLabelStyle.Render("Example:"),
				descStyle.Render(cmd.example)))
		}
	}

	s.WriteString("\n")
	kbHeader := headerStyle.Render(" KEYBOARD SHORTCUTS ")
	s.WriteString(kbHeader)
	s.WriteString("\n\n")

	// Keyboard shortcuts with style
	keyStyle := lipgloss.NewStyle().Foreground(accentGreen).Bold(true)
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("↑/k      "), descStyle.Render("Move cursor up")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("↓/j      "), descStyle.Render("Move cursor down")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("Enter    "), descStyle.Render("Execute command")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("Esc      "), descStyle.Render("Clear input")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("?        "), descStyle.Render("Toggle help")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("Ctrl+C   "), descStyle.Render("Quit")))

	s.WriteString("\n")

	// Scroll help
	s.WriteString(headerStyle.Render(" SCROLL NAVIGATION "))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("↑/k      "), descStyle.Render("Scroll up one line")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("↓/j      "), descStyle.Render("Scroll down one line")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("Space/f  "), descStyle.Render("Page down")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("b        "), descStyle.Render("Page up")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("g        "), descStyle.Render("Go to top")))
	s.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render("G        "), descStyle.Render("Go to bottom")))

	s.WriteString("\n")
	footerStyle := lipgloss.NewStyle().Foreground(fgDim).Italic(true)
	s.WriteString(footerStyle.Render("  Press 'q', 'Esc', or 'Ctrl+C' to return to the main view  "))
	s.WriteString("\n")

	return s.String()
}

// renderHelpView renders the help screen using viewport
func (m Model) renderHelpView() string {
	return m.viewport.View()
}

// renderBannerView renders the startup banner screen
func (m Model) renderBannerView() string {
	// Left column: Banner
	var leftCol strings.Builder
	leftCol.WriteString("\n\n")

	// Render the ASCII art banner
	banner := GetBanner()
	leftCol.WriteString(bannerStyle.Render(banner))
	leftCol.WriteString("\n\n")

	// Render subtitle
	subtitle := GetSubtitle()
	leftCol.WriteString(bannerSubtitleStyle.Render(subtitle))
	leftCol.WriteString("\n\n")

	// Render version
	version := GetVersion()
	leftCol.WriteString(bannerVersionStyle.Render(version))
	leftCol.WriteString("\n")

	// Render "press any key" prompt
	leftCol.WriteString(bannerPromptStyle.Render("Press any key to continue..."))

	// Right column: Recent Todos
	rightCol := m.renderBannerTodoBox()

	// Join columns horizontally with spacing
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftCol.String(),
		strings.Repeat(" ", 4), // spacing between columns
		rightCol,
	)

	// Add top padding
	return "\n\n" + content
}

// renderBannerTodoBox renders the todo list box on the banner screen
func (m Model) renderBannerTodoBox() string {
	var content strings.Builder

	// Title
	content.WriteString(bannerTodoTitleStyle.Render("📋 Recent Todos"))
	content.WriteString("\n\n")

	// Get oldest 5 todos (sorted by creation date)
	if len(m.todos) == 0 {
		content.WriteString(emptyStyle.Render("No todos yet!"))
	} else {
		// Display up to 5 oldest todos
		count := len(m.todos)
		if count > 5 {
			count = 5
		}

		for i := 0; i < count; i++ {
			todo := m.todos[i]

			// Number in green
			number := bannerTodoNumberStyle.Render(fmt.Sprintf("%d.", i+1))

			// Title - truncate if longer than 15 characters
			title := truncateTodoTitle(todo.Title, 15)

			// Format: number title
			line := fmt.Sprintf("%s %s", number, bannerTodoItemStyle.Render(title))

			content.WriteString(line)
			content.WriteString("\n")
		}

		// Show count if there are more
		if len(m.todos) > 5 {
			content.WriteString("\n")
			content.WriteString(helpStyle.Render(fmt.Sprintf("+ %d more todos...", len(m.todos)-5)))
		}
	}

	// Wrap in border box
	return bannerTodoBoxStyle.Render(content.String())
}

// truncateTodoTitle truncates a todo title to maxLength characters, adding "..." if truncated
func truncateTodoTitle(title string, maxLength int) string {
	// Convert to runes to handle multibyte characters correctly
	runes := []rune(title)
	if len(runes) <= maxLength {
		return title
	}
	return string(runes[:maxLength]) + "..."
}

// renderAddTodoView renders the add todo screen
func (m Model) renderAddTodoView() string {
	var s strings.Builder

	// Title with dark background
	s.WriteString(titleStyle.Render(" ➕ Add New Todo "))
	s.WriteString("\n\n")

	// Show current step
	stepIndicator := ""
	if m.addTodoStep == 0 {
		stepIndicator = headerStyle.Render(" Step 1/2: Enter Title ")
	} else {
		stepIndicator = headerStyle.Render(" Step 2/2: Enter Description (Optional) ")
	}
	s.WriteString(stepIndicator)
	s.WriteString("\n\n")

	// Show previously entered title if on step 2
	if m.addTodoStep == 1 {
		s.WriteString(messageStyle.Render(fmt.Sprintf("Title: %s", m.addTodoTitle)))
		s.WriteString("\n\n")
	}

	// Input field
	s.WriteString(m.input.View())
	s.WriteString("\n")

	// Error message if any
	if m.err != nil {
		s.WriteString("\n")
		s.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		s.WriteString("\n")
	}

	// Help text
	s.WriteString("\n")
	if m.addTodoStep == 0 {
		s.WriteString(helpStyle.Render("Press Enter to continue | Esc to cancel"))
	} else {
		s.WriteString(helpStyle.Render("Press Enter to save | Esc to go back"))
	}

	return s.String()
}

// renderEditTodoView renders the edit todo screen
func (m Model) renderEditTodoView() string {
	var s strings.Builder

	// Title with dark background
	s.WriteString(titleStyle.Render(fmt.Sprintf(" ✏️  Edit Todo #%d ", m.editTodoID)))
	s.WriteString("\n\n")

	// Show current step
	stepIndicator := ""
	if m.editTodoStep == 0 {
		stepIndicator = headerStyle.Render(" Step 1/2: Edit Title ")
	} else {
		stepIndicator = headerStyle.Render(" Step 2/2: Edit Description (Optional) ")
	}
	s.WriteString(stepIndicator)
	s.WriteString("\n\n")

	// Show previously entered title if on step 2
	if m.editTodoStep == 1 {
		s.WriteString(messageStyle.Render(fmt.Sprintf("Title: %s", m.editTodoTitle)))
		s.WriteString("\n\n")
	}

	// Input field
	s.WriteString(m.input.View())
	s.WriteString("\n")

	// Error message if any
	if m.err != nil {
		s.WriteString("\n")
		s.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		s.WriteString("\n")
	}

	// Help text
	s.WriteString("\n")
	if m.editTodoStep == 0 {
		s.WriteString(helpStyle.Render("Press Enter to continue | Esc to cancel"))
	} else {
		s.WriteString(helpStyle.Render("Press Enter to save | Esc to go back"))
	}

	return s.String()
}

// renderPomodoroView renders the Pomodoro timer screen
func (m Model) renderPomodoroView() string {
	var s strings.Builder

	// Title with dark background
	s.WriteString(titleStyle.Render(" 🍅 Pomodoro Timer "))
	s.WriteString("\n\n")

	// Timer display - large and centered
	minutes := m.pomoSecondsLeft / 60
	seconds := m.pomoSecondsLeft % 60
	timerText := fmt.Sprintf("%02d:%02d", minutes, seconds)

	// Make the timer text large by adding spacing
	largeTimerText := ""
	for _, char := range timerText {
		largeTimerText += string(char) + " "
	}
	largeTimerDisplay := lipgloss.NewStyle().
		Foreground(accentGreen).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		MarginTop(3).
		MarginBottom(3).
		Render(largeTimerText)

	s.WriteString(largeTimerDisplay)
	s.WriteString("\n")

	// Show which todo is being worked on
	if m.pomoTodoID > 0 {
		// Find the todo
		var todoTitle string
		for _, todo := range m.todos {
			if todo.ID == m.pomoTodoID {
				todoTitle = todo.Title
				break
			}
		}

		if todoTitle != "" {
			taskInfo := fmt.Sprintf("Working on: #%d - %s", m.pomoTodoID, todoTitle)
			s.WriteString(messageStyle.Render(taskInfo))
			s.WriteString("\n\n")
		}
	} else {
		s.WriteString(helpStyle.Render("General Pomodoro session"))
		s.WriteString("\n\n")
	}

	// Status indicator
	if m.pomoRunning {
		s.WriteString(messageStyle.Render("⏱️  Timer is running..."))
	} else {
		s.WriteString(helpStyle.Render("Timer paused"))
	}
	s.WriteString("\n\n")

	// Help text
	s.WriteString(helpStyle.Render("Press Esc or Enter to stop timer"))

	return s.String()
}
