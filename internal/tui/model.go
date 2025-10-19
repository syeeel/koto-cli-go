package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/syeeel/koto-cli-go/internal/model"
	"github.com/syeeel/koto-cli-go/internal/service"
)

// ViewMode represents the current view mode of the TUI
type ViewMode int

const (
	// ViewModeBanner shows the startup banner
	ViewModeBanner ViewMode = iota
	// ViewModeList shows the list of todos
	ViewModeList
	// ViewModeHelp shows the help screen
	ViewModeHelp
)

// Model represents the Bubbletea model for the TUI
type Model struct {
	service  *service.TodoService
	todos    []*model.Todo
	cursor   int
	viewMode ViewMode
	input    textinput.Model
	message  string
	err      error
	width    int
	height   int
	quitting bool
}

// NewModel creates a new TUI model
func NewModel(service *service.TodoService) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter command (type /help for help)"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80

	return Model{
		service:  service,
		todos:    []*model.Todo{},
		cursor:   0,
		viewMode: ViewModeBanner,
		input:    ti,
		quitting: false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		loadTodos(m.service),
	)
}
