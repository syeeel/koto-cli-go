package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

	// Title
	s.WriteString(titleStyle.Render("📝 koto - ToDo Manager"))
	s.WriteString("\n\n")

	// Todo list
	if len(m.todos) == 0 {
		s.WriteString(emptyStyle.Render("No todos yet. Use /add to create your first todo!"))
		s.WriteString("\n")
	} else {
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
	s.WriteString(helpStyle.Render("Commands: /add, /list, /done, /delete, /edit, /help | Navigate: ↑/↓ or j/k | Help: ? | Quit: Ctrl+C"))

	return s.String()
}

// renderTodoItem renders a single todo item
func (m Model) renderTodoItem(index int, todo *model.Todo) string {
	var s strings.Builder

	// Cursor
	cursor := "  "
	if m.cursor == index {
		cursor = selectedStyle.Render("> ")
	} else {
		cursor = "  "
	}
	s.WriteString(cursor)

	// Status checkbox
	status := "⬜"
	if todo.IsCompleted() {
		status = "✅"
	}
	s.WriteString(status)
	s.WriteString(" ")

	// Priority indicator
	priority := m.renderPriority(todo.Priority)
	s.WriteString(priority)
	s.WriteString(" ")

	// ID
	s.WriteString(fmt.Sprintf("[%d] ", todo.ID))

	// Title (with strike-through for completed)
	title := todo.Title
	if todo.IsCompleted() {
		title = completedItemStyle.Render(title)
	} else {
		if m.cursor == index {
			title = selectedStyle.Render(title)
		} else {
			title = todoItemStyle.Render(title)
		}
	}
	s.WriteString(title)

	// Due date indicator (if overdue)
	if todo.IsOverdue() {
		s.WriteString(" ")
		s.WriteString(errorStyle.Render("⚠ OVERDUE"))
	}

	return s.String()
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

	s.WriteString(titleStyle.Render("📖 koto - Help"))
	s.WriteString("\n\n")

	s.WriteString("COMMANDS:\n\n")

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
		{"/quit", "Quit the application", "/quit"},
	}

	for _, cmd := range commands {
		if cmd.command == "" {
			s.WriteString("\n")
			continue
		}
		s.WriteString(fmt.Sprintf("  %-45s %s\n", selectedStyle.Render(cmd.command), cmd.desc))
		if cmd.example != "" {
			s.WriteString(fmt.Sprintf("    %s %s\n", helpStyle.Render("Example:"), cmd.example))
		}
	}

	s.WriteString("\n\nKEYBOARD SHORTCUTS:\n\n")
	s.WriteString("  ↑/k       Move cursor up\n")
	s.WriteString("  ↓/j       Move cursor down\n")
	s.WriteString("  Enter     Execute command\n")
	s.WriteString("  Esc       Clear input\n")
	s.WriteString("  ?         Toggle help\n")
	s.WriteString("  Ctrl+C    Quit\n")

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("Press 'q', 'Esc', or 'Ctrl+C' to return to the main view"))

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
