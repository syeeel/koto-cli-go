package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
			}
			return m, nil
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

	return m, parseAndExecuteCommand(m.service, value)
}
