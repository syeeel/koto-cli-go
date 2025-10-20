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
		// Header with fixed widths: No.(4), Title(40), Description(40), Create Date(11)
		headerNo := padStringToWidth("No.", 4)
		headerTitle := padStringToWidth("Title", 40)
		headerDesc := padStringToWidth("Description", 40)
		headerDate := padStringToWidth("Create Date", 11)
		header := fmt.Sprintf(" %s   %s   %s   %s ", headerNo, headerTitle, headerDesc, headerDate)
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
	s.WriteString(helpStyle.Render("Commands: /add, /list, /done, /delete, /edit, /help | Navigate: ↑/↓ or j/k | Help: ? | Quit: /exit or Ctrl+C"))

	return s.String()
}

// renderTodoItem renders a single todo item in table format
func (m Model) renderTodoItem(index int, todo *model.Todo) string {
	// No. (ID) - width: 4
	no := fmt.Sprintf("%d", todo.ID)
	no = padStringToWidth(no, 4)

	// Title - width: 40 (display width, not character count)
	title := truncateStringByWidth(todo.Title, 40)
	title = padStringToWidth(title, 40)

	// Description - width: 40 (display width, not character count)
	desc := truncateStringByWidth(todo.Description, 40)
	desc = padStringToWidth(desc, 40)

	// Create Date (format: YYYY-MM-DD) - width: 11
	createDate := todo.CreatedAt.Format("2006-01-02")
	createDate = padStringToWidth(createDate, 11)

	// Build the row with spacing (no vertical separators for cleaner look)
	row := fmt.Sprintf(" %s   %s   %s   %s ", no, title, desc, createDate)

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

// renderHelpView renders the help screen
func (m Model) renderHelpView() string {
	var s strings.Builder

	// Title with dark background
	title := titleStyle.Render(" 📖 koto - Help ")
	s.WriteString(title)
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
		{"/add <title>", "Add a new todo", "/add Buy groceries"},
		{"/add <title> --desc=<description>", "Add todo with description", "/add Study --desc=\"Chapter 5\""},
		{"/add <title> --priority=<low|medium|high>", "Add todo with priority", "/add Report --priority=high"},
		{"/add <title> --due=<YYYY-MM-DD>", "Add todo with due date", "/add Project --due=2025-12-31"},
		{"", "", ""},
		{"/list", "List all todos", "/list"},
		{"/list --status=<pending|completed>", "List by status", "/list --status=pending"},
		{"", "", ""},
		{"/done <id>", "Mark todo as completed", "/done 1"},
		{"/delete <id>", "Delete a todo", "/delete 2"},
		{"/edit <id> --title=<new title>", "Edit todo title", "/edit 1 --title=\"New title\""},
		{"/edit <id> --priority=<low|medium|high>", "Edit todo priority", "/edit 1 --priority=high"},
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
	footerStyle := lipgloss.NewStyle().Foreground(fgDim).Italic(true)
	s.WriteString(footerStyle.Render("  Press 'q', 'Esc', or 'Ctrl+C' to return to the main view  "))
	s.WriteString("\n")

	return s.String()
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
