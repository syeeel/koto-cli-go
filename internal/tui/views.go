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
	case ViewModeDetail:
		return m.renderDetailView()
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
		// Header with fixed widths: No.(5), Title(60), Priority(12), Total time(12), Create Date(13)
		headerNo := padStringToWidth("No.", 5)
		headerTitle := padStringToWidth("Title", 60)
		headerPriority := padStringToWidth("Priority", 12)
		headerTime := padStringToWidth("Total time", 12)
		headerDate := padStringToWidth("Create Date", 13)
		header := fmt.Sprintf(" %s   %s   %s   %s   %s ", headerNo, headerTitle, headerPriority, headerTime, headerDate)
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
	s.WriteString(helpStyle.Render("Commands: /add, /list, /done, /edit, /pomo, /help | Navigate: ↑/↓ or j/k | Help: ? | Quit: /exit or Ctrl+C"))

	return s.String()
}

// renderTodoItem renders a single todo item in table format
func (m Model) renderTodoItem(index int, todo *model.Todo) string {
	// No. (ID) - width: 5
	no := fmt.Sprintf("%d", todo.ID)
	no = padStringToWidth(no, 5)

	// Title - width: 60 (display width, not character count)
	title := truncateStringByWidth(todo.Title, 60)
	title = padStringToWidth(title, 60)

	// Priority - width: 12 (display width, not character count)
	// First create the plain text, pad it, then apply styling
	var priorityText string
	var priorityColor lipgloss.Color
	switch todo.Priority {
	case model.PriorityHigh:
		priorityText = "High"
		priorityColor = lipgloss.Color("196")
	case model.PriorityMedium:
		priorityText = "Medium"
		priorityColor = lipgloss.Color("220")
	case model.PriorityLow:
		priorityText = "Low"
		priorityColor = lipgloss.Color("82")
	}
	// Pad the plain text first
	priorityPadded := padStringToWidth(priorityText, 12)

	// For styling, check if this item is selected
	var priority string
	if m.cursor == index {
		// When selected, don't apply color styling (let selectedStyle handle it)
		priority = priorityPadded
	} else {
		// When not selected, apply color styling
		priority = lipgloss.NewStyle().
			Foreground(priorityColor).
			Render(priorityPadded)
	}

	// Total time - width: 12
	totalTime := todo.GetWorkDurationFormatted()
	if totalTime == "" {
		totalTime = "-"
	}
	totalTime = padStringToWidth(totalTime, 12)

	// Create Date (format: YYYY-MM-DD) - width: 13
	createDate := todo.CreatedAt.Format("2006-01-02")
	createDate = padStringToWidth(createDate, 13)

	// Build the row with spacing (no vertical separators for cleaner look)
	row := fmt.Sprintf(" %s   %s   %s   %s   %s ", no, title, priority, totalTime, createDate)

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
		{"", "  → Step 3: Select priority (1-3)", ""},
		{"", "", ""},
		{"/list", "List all todos", "/list"},
		{"/list --status=<pending|completed>", "List by status", "/list --status=pending"},
		{"", "", ""},
		{"/done <id>", "Delete a todo", "/done 1"},
		{"/edit <id>", "Edit a todo (interactive)", "/edit 1"},
		{"", "  → Step 1: Edit title", ""},
		{"", "  → Step 2: Edit description (optional)", ""},
		{"", "  → Step 3: Select priority (1-3)", ""},
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
		stepIndicator = headerStyle.Render(" Step 1/3: Enter Title ")
	} else if m.addTodoStep == 1 {
		stepIndicator = headerStyle.Render(" Step 2/3: Enter Description (Optional) ")
	} else {
		stepIndicator = headerStyle.Render(" Step 3/3: Select Priority ")
	}
	s.WriteString(stepIndicator)
	s.WriteString("\n\n")

	// Show previously entered data
	if m.addTodoStep >= 1 {
		s.WriteString(messageStyle.Render(fmt.Sprintf("Title: %s", m.addTodoTitle)))
		s.WriteString("\n\n")
	}
	if m.addTodoStep == 2 {
		if m.addTodoDescription != "" {
			s.WriteString(messageStyle.Render(fmt.Sprintf("Description: %s", m.addTodoDescription)))
		} else {
			s.WriteString(emptyStyle.Render("Description: (none)"))
		}
		s.WriteString("\n\n")
		// Show priority options
		s.WriteString(lipgloss.NewStyle().Foreground(accentGreen).Render("Priority Options:"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("  1 = Low"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("  2 = Medium (default)"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("  3 = High"))
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
		stepIndicator = headerStyle.Render(" Step 1/3: Edit Title ")
	} else if m.editTodoStep == 1 {
		stepIndicator = headerStyle.Render(" Step 2/3: Edit Description (Optional) ")
	} else {
		stepIndicator = headerStyle.Render(" Step 3/3: Select Priority ")
	}
	s.WriteString(stepIndicator)
	s.WriteString("\n\n")

	// Show previously entered data
	if m.editTodoStep >= 1 {
		s.WriteString(messageStyle.Render(fmt.Sprintf("Title: %s", m.editTodoTitle)))
		s.WriteString("\n\n")
	}
	if m.editTodoStep == 2 {
		if m.editTodoDescription != "" {
			s.WriteString(messageStyle.Render(fmt.Sprintf("Description: %s", m.editTodoDescription)))
		} else {
			s.WriteString(emptyStyle.Render("Description: (none)"))
		}
		s.WriteString("\n\n")
		// Show priority options
		s.WriteString(lipgloss.NewStyle().Foreground(accentGreen).Render("Priority Options:"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("  1 = Low"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("  2 = Medium (default)"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("  3 = High"))
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

	// Title with tomato emoji
	s.WriteString(titleStyle.Render(" 🍅 Pomodoro Timer "))
	s.WriteString("\n\n")

	// Calculate progress
	totalSeconds := 1500 // 25 minutes
	elapsedSeconds := totalSeconds - m.pomoSecondsLeft
	progressPercent := float64(elapsedSeconds) / float64(totalSeconds)

	// Timer display - large numbers with gradient
	minutes := m.pomoSecondsLeft / 60
	seconds := m.pomoSecondsLeft % 60

	// Render large timer numbers with pink gradient
	largeTimer := m.renderLargeTimer(minutes, seconds)
	s.WriteString(lipgloss.PlaceHorizontal(100, lipgloss.Center, largeTimer))
	s.WriteString("\n\n")

	// Progress bar with pink gradient
	progressView := m.renderPinkProgressBar(progressPercent)
	s.WriteString(lipgloss.PlaceHorizontal(100, lipgloss.Center, progressView))
	s.WriteString("\n\n")

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
			taskInfoBox := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#585b70")).
				Padding(1, 2).
				Width(70).
				Render(
					lipgloss.NewStyle().
						Foreground(lipgloss.Color("213")).
						Bold(true).
						Render(fmt.Sprintf("📋 Task #%d", m.pomoTodoID)) +
						"\n" +
						lipgloss.NewStyle().
							Foreground(fgDefault).
							Render(todoTitle),
				)
			s.WriteString(lipgloss.PlaceHorizontal(100, lipgloss.Center, taskInfoBox))
			s.WriteString("\n\n")
		}
	} else {
		infoBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#585b70")).
			Padding(1, 2).
			Width(70).
			Align(lipgloss.Center).
			Render(
				lipgloss.NewStyle().
					Foreground(fgDim).
					Italic(true).
					Render("General Pomodoro session"),
			)
		s.WriteString(lipgloss.PlaceHorizontal(100, lipgloss.Center, infoBox))
		s.WriteString("\n\n")
	}

	// Status indicator
	statusText := ""
	if m.pomoCompleted {
		// Timer completed - show alarm message
		statusText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Render("🔔 Timer Complete! Press Enter or Esc to stop alarm")
	} else if m.pomoRunning {
		statusText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(true).
			Render("⏱️  Timer is running...")
	} else {
		statusText = lipgloss.NewStyle().
			Foreground(fgDim).
			Render("Timer paused")
	}
	s.WriteString(lipgloss.PlaceHorizontal(100, lipgloss.Center, statusText))
	s.WriteString("\n\n")

	// Help text in a subtle box
	helpBox := lipgloss.NewStyle().
		Foreground(fgDim).
		Italic(true).
		Render("Press Esc or Enter to stop timer")
	s.WriteString(lipgloss.PlaceHorizontal(100, lipgloss.Center, helpBox))

	return s.String()
}

// renderLargeTimer renders large timer numbers with pink gradient
func (m Model) renderLargeTimer(minutes, seconds int) string {
	// ASCII art style large numbers (simplified)
	minutesStr := fmt.Sprintf("%02d", minutes)
	secondsStr := fmt.Sprintf("%02d", seconds)

	// Create gradient using banner pink color (213) and variations
	pink1 := lipgloss.Color("213") // Banner pink (base)
	pink2 := lipgloss.Color("212") // Slightly lighter
	pink3 := lipgloss.Color("205") // Brighter variant

	// Render each digit with spacing and gradient
	var lines [5]string
	digits := minutesStr + ":" + secondsStr

	for i, char := range digits {
		var color lipgloss.Color
		if i < 2 {
			color = pink1
		} else if i == 2 {
			color = pink2
		} else {
			color = pink3
		}

		digitArt := getDigitArt(char)
		for lineIdx, line := range digitArt {
			styledLine := lipgloss.NewStyle().
				Foreground(color).
				Bold(true).
				Render(line)
			if i > 0 {
				lines[lineIdx] += "  " // Spacing between digits
			}
			lines[lineIdx] += styledLine
		}
	}

	// Join all lines
	var result strings.Builder
	for _, line := range lines {
		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}

// interpolateColor interpolates between two RGB colors
func interpolateColor(startR, startG, startB, endR, endG, endB int, ratio float64) string {
	r := int(float64(startR) + (float64(endR)-float64(startR))*ratio)
	g := int(float64(startG) + (float64(endG)-float64(startG))*ratio)
	b := int(float64(startB) + (float64(endB)-float64(startB))*ratio)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// renderPinkProgressBar renders a beautiful progress bar with smooth 3-color gradient (pink→blue→green)
func (m Model) renderPinkProgressBar(percent float64) string {
	width := 60
	filled := int(percent * float64(width))

	// Three-color gradient: Pink → Blue → Green
	// Start: Light pink (#ff69b4)
	pinkR, pinkG, pinkB := 255, 105, 180
	// Middle: Royal blue (#4169e1)
	blueR, blueG, blueB := 65, 105, 225
	// End: Primary green (#06c775)
	greenR, greenG, greenB := 6, 199, 117

	var bar strings.Builder

	// Filled portion with smooth 3-color gradient
	for i := 0; i < filled; i++ {
		var color string

		// Calculate position ratio (0.0 to 1.0)
		positionRatio := float64(i) / float64(width-1)

		if positionRatio <= 0.5 {
			// First half: Pink → Blue
			segmentRatio := positionRatio / 0.5 // 0.0 to 1.0 within this segment
			color = interpolateColor(pinkR, pinkG, pinkB, blueR, blueG, blueB, segmentRatio)
		} else {
			// Second half: Blue → Green
			segmentRatio := (positionRatio - 0.5) / 0.5 // 0.0 to 1.0 within this segment
			color = interpolateColor(blueR, blueG, blueB, greenR, greenG, greenB, segmentRatio)
		}

		styled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render("█")
		bar.WriteString(styled)
	}

	// Empty portion
	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#313244"))
	for i := filled; i < width; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}

	// Add percentage text in primary green
	percentText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06c775")).
		Bold(true).
		Render(fmt.Sprintf(" %3.0f%%", percent*100))
	bar.WriteString(percentText)

	return bar.String()
}

// getDigitArt returns ASCII art for a single digit or colon
func getDigitArt(char rune) [5]string {
	switch char {
	case '0':
		return [5]string{
			"█████",
			"█   █",
			"█   █",
			"█   █",
			"█████",
		}
	case '1':
		return [5]string{
			"  ██ ",
			" ███ ",
			"  ██ ",
			"  ██ ",
			"█████",
		}
	case '2':
		return [5]string{
			"█████",
			"    █",
			"█████",
			"█    ",
			"█████",
		}
	case '3':
		return [5]string{
			"█████",
			"    █",
			"█████",
			"    █",
			"█████",
		}
	case '4':
		return [5]string{
			"█   █",
			"█   █",
			"█████",
			"    █",
			"    █",
		}
	case '5':
		return [5]string{
			"█████",
			"█    ",
			"█████",
			"    █",
			"█████",
		}
	case '6':
		return [5]string{
			"█████",
			"█    ",
			"█████",
			"█   █",
			"█████",
		}
	case '7':
		return [5]string{
			"█████",
			"    █",
			"   ██",
			"  ██ ",
			" ██  ",
		}
	case '8':
		return [5]string{
			"█████",
			"█   █",
			"█████",
			"█   █",
			"█████",
		}
	case '9':
		return [5]string{
			"█████",
			"█   █",
			"█████",
			"    █",
			"█████",
		}
	case ':':
		return [5]string{
			"     ",
			" ██  ",
			"     ",
			" ██  ",
			"     ",
		}
	default:
		return [5]string{"", "", "", "", ""}
	}
}

// renderDetailView renders the todo detail screen
func (m Model) renderDetailView() string {
	var s strings.Builder

	// Find the todo to display
	var targetTodo *model.Todo
	for _, todo := range m.todos {
		if todo.ID == m.detailTodoID {
			targetTodo = todo
			break
		}
	}

	if targetTodo == nil {
		errorBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentRed).
			Padding(1, 2).
			Width(60).
			Render(errorStyle.Render("Todo not found"))
		s.WriteString(errorBox)
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Esc to return to main view"))
		return s.String()
	}

	// Main title
	s.WriteString(titleStyle.Render(fmt.Sprintf(" 📋 Todo Details #%d ", targetTodo.ID)))
	s.WriteString("\n\n")

	// Title field box
	titleLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Title")
	titleContent := lipgloss.NewStyle().
		Foreground(fgDefault).
		Render(targetTodo.Title)
	titleBoxContent := titleLabel + "\n" + titleContent
	titleBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585b70")).
		Padding(1, 2).
		Width(110)
	titleBox := titleBoxStyle.Render(titleBoxContent)
	s.WriteString(titleBox)
	s.WriteString("\n\n")

	// Description field box (with increased height)
	descLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Description")
	var descContent string
	if targetTodo.Description != "" {
		descContent = lipgloss.NewStyle().
			Foreground(fgDefault).
			Render(targetTodo.Description)
	} else {
		descContent = emptyStyle.Render("(no description)")
	}
	descBoxContent := descLabel + "\n" + descContent
	descBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585b70")).
		Padding(2, 2).  // Increased vertical padding
		Height(8).      // Set explicit height for 3x size
		Width(110)
	descBox := descBoxStyle.Render(descBoxContent)
	s.WriteString(descBox)
	s.WriteString("\n\n")

	// Priority, Total Work Time, and Timestamps (3 columns in one row)

	// Priority content
	var priorityValue string
	switch targetTodo.Priority {
	case model.PriorityHigh:
		priorityValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render("High")
	case model.PriorityMedium:
		priorityValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Render("Medium")
	case model.PriorityLow:
		priorityValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Render("Low")
	}

	// Work duration content
	var workValue string
	if targetTodo.WorkDuration > 0 {
		hours := targetTodo.WorkDuration / 60
		minutes := targetTodo.WorkDuration % 60
		workValue = lipgloss.NewStyle().
			Foreground(fgDefault).
			Render(fmt.Sprintf("🍅 %dh %dm", hours, minutes))
	} else {
		workValue = emptyStyle.Render("(no records)")
	}

	// Timestamps content
	createdStr := lipgloss.NewStyle().
		Foreground(fgDim).
		Render("Created: ") +
		lipgloss.NewStyle().
			Foreground(fgDefault).
			Render(targetTodo.CreatedAt.Format("2006-01-02 15:04:05"))
	timestampContent := createdStr

	// Priority box
	priorityLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Priority")
	priorityBoxContent := priorityLabel + "\n" + priorityValue
	priorityBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585b70")).
		Padding(1, 2).
		Width(34)
	priorityBox := priorityBoxStyle.Render(priorityBoxContent)

	// Work time box
	workLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Total Work Time")
	workBoxContent := workLabel + "\n" + workValue
	workBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585b70")).
		Padding(1, 2).
		Width(34)
	workBox := workBoxStyle.Render(workBoxContent)

	// Created timestamp box
	timestampLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Created")
	timestampBoxContent := timestampLabel + "\n" + timestampContent
	timestampBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585b70")).
		Padding(1, 2).
		Width(38)
	timestampBox := timestampBoxStyle.Render(timestampBoxContent)

	// Join three boxes horizontally
	threeColumnRow := lipgloss.JoinHorizontal(lipgloss.Top, priorityBox, "  ", workBox, "  ", timestampBox)
	s.WriteString(threeColumnRow)
	s.WriteString("\n\n")

	// Help text
	s.WriteString(helpStyle.Render("Press Enter to return | e to edit | d to done"))

	return s.String()
}
