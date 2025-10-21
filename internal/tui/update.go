package tui

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/syeeel/koto-cli-go/internal/model"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update viewport size if in help mode
		if m.viewMode == ViewModeHelp {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 2
		}
		return m, nil

	case tea.KeyMsg:
		// Handle banner view - any key transitions to list view
		if m.viewMode == ViewModeBanner {
			m.viewMode = ViewModeList
			return m, nil
		}

		// Handle view mode specific keys
		if m.viewMode == ViewModeHelp {
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				m.viewMode = ViewModeList
				return m, nil
			case "up", "k":
				m.viewport.LineUp(1)
				return m, nil
			case "down", "j":
				m.viewport.LineDown(1)
				return m, nil
			case "pgup", "b":
				m.viewport.ViewUp()
				return m, nil
			case "pgdown", "f", " ":
				m.viewport.ViewDown()
				return m, nil
			case "g":
				m.viewport.GotoTop()
				return m, nil
			case "G":
				m.viewport.GotoBottom()
				return m, nil
			}
			return m, nil
		}

		// Handle add todo view
		if m.viewMode == ViewModeAddTodo {
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "esc":
				// Cancel and return to list view, or go back to previous step
				if m.addTodoStep == 0 {
					// Cancel the whole operation
					m.viewMode = ViewModeList
					m.input.Placeholder = "Enter command (type /help for help)"
					m.input.SetValue("")
					m.message = "Add todo cancelled"
				} else {
					// Go back to previous step
					m.addTodoStep = 0
					m.input.Placeholder = "Enter todo title..."
					m.input.SetValue(m.addTodoTitle)
				}
				return m, nil

			case "enter":
				return m.handleAddTodoEnter()
			}

			// Update input for add todo view
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		// Handle edit todo view
		if m.viewMode == ViewModeEditTodo {
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "esc":
				// Cancel and return to list view, or go back to previous step
				if m.editTodoStep == 0 {
					// Cancel the whole operation
					m.viewMode = ViewModeList
					m.input.Placeholder = "Enter command (type /help for help)"
					m.input.SetValue("")
					m.message = "Edit cancelled"
				} else {
					// Go back to previous step
					m.editTodoStep = 0
					m.input.Placeholder = "Edit todo title..."
					m.input.SetValue(m.editTodoTitle)
				}
				return m, nil

			case "enter":
				return m.handleEditTodoEnter()
			}

			// Update input for edit todo view
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		// Handle list view keys
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			// Clear input and message
			m.input.SetValue("")
			m.message = ""
			m.err = nil
			return m, nil

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
			return m, nil

		case "?":
			m.viewMode = ViewModeHelp
			// Initialize viewport for help view
			m.viewport = viewport.New(m.width, m.height-2)
			m.viewport.SetContent(m.renderHelpContent())
			return m, nil
		}

	case todosLoadedMsg:
		m.todos = msg.todos
		m.err = msg.err
		// Adjust cursor if it's out of bounds
		if m.cursor >= len(m.todos) && len(m.todos) > 0 {
			m.cursor = len(m.todos) - 1
		}
		if len(m.todos) == 0 {
			m.cursor = 0
		}
		return m, nil

	case commandExecutedMsg:
		m.message = msg.message
		m.err = msg.err
		// Reload todos after command execution
		return m, loadTodos(m.service)
	}

	// Update the text input
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// handleEnter processes the enter key press
func (m *Model) handleEnter() (tea.Model, tea.Cmd) {
	value := m.input.Value()
	m.input.SetValue("")
	m.message = ""
	m.err = nil

	// If input is empty, do nothing
	if value == "" {
		return m, nil
	}

	// Check if command is /add - switch to add todo view
	if value == "/add" {
		m.viewMode = ViewModeAddTodo
		m.addTodoStep = 0
		m.addTodoTitle = ""
		m.addTodoDescription = ""
		m.input.Placeholder = "Enter todo title..."
		m.input.SetValue("")
		return m, nil
	}

	// Check if command is /edit <id> - switch to edit todo view
	if strings.HasPrefix(value, "/edit ") {
		parts := strings.Fields(value)
		if len(parts) == 2 {
			id, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				m.err = errors.New("invalid todo ID")
				return m, nil
			}

			// Find the todo with the given ID
			var targetTodo *model.Todo
			for _, todo := range m.todos {
				if todo.ID == id {
					targetTodo = todo
					break
				}
			}

			if targetTodo == nil {
				m.err = errors.New("todo not found")
				return m, nil
			}

			// Switch to edit mode with existing data
			m.viewMode = ViewModeEditTodo
			m.editTodoID = id
			m.editTodoTitle = targetTodo.Title
			m.editTodoDescription = targetTodo.Description
			m.editTodoStep = 0
			m.input.Placeholder = "Edit todo title..."
			m.input.SetValue(targetTodo.Title)
			m.err = nil
			return m, nil
		}
	}

	return m, parseAndExecuteCommand(m.service, value)
}

// handleAddTodoEnter processes the enter key press in add todo view
func (m *Model) handleAddTodoEnter() (tea.Model, tea.Cmd) {
	value := m.input.Value()

	switch m.addTodoStep {
	case 0: // Title input
		if value == "" {
			m.err = errors.New("title cannot be empty")
			return m, nil
		}
		// Save title and move to description step
		m.addTodoTitle = value
		m.addTodoStep = 1
		m.input.SetValue("")
		m.input.Placeholder = "Enter description (optional, press Enter to skip)..."
		m.err = nil
		return m, nil

	case 1: // Description input
		// Save description
		m.addTodoDescription = value

		// Create the todo
		ctx := context.Background()
		_, err := m.service.AddTodo(ctx, m.addTodoTitle, m.addTodoDescription, model.PriorityMedium, nil)
		if err != nil {
			m.err = err
			return m, nil
		}

		// Reset and return to list view
		m.viewMode = ViewModeList
		m.input.Placeholder = "Enter command (type /help for help)"
		m.input.SetValue("")
		m.message = "Todo added successfully"
		m.err = nil

		// Reload todos
		return m, loadTodos(m.service)
	}

	return m, nil
}

// handleEditTodoEnter processes the enter key press in edit todo view
func (m *Model) handleEditTodoEnter() (tea.Model, tea.Cmd) {
	value := m.input.Value()

	switch m.editTodoStep {
	case 0: // Title input
		if value == "" {
			m.err = errors.New("title cannot be empty")
			return m, nil
		}
		// Save title and move to description step
		m.editTodoTitle = value
		m.editTodoStep = 1
		m.input.SetValue(m.editTodoDescription)
		m.input.Placeholder = "Edit description (optional, press Enter to save)..."
		m.err = nil
		return m, nil

	case 1: // Description input
		// Save description
		m.editTodoDescription = value

		// Update the todo
		ctx := context.Background()
		err := m.service.EditTodo(ctx, m.editTodoID, m.editTodoTitle, m.editTodoDescription, model.PriorityMedium, nil)
		if err != nil {
			m.err = err
			return m, nil
		}

		// Reset and return to list view
		m.viewMode = ViewModeList
		m.input.Placeholder = "Enter command (type /help for help)"
		m.input.SetValue("")
		m.message = "Todo updated successfully"
		m.err = nil

		// Reload todos
		return m, loadTodos(m.service)
	}

	return m, nil
}
