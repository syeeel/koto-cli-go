package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	// ViewModeAddTodo shows the add todo screen
	ViewModeAddTodo
	// ViewModeEditTodo shows the edit todo screen
	ViewModeEditTodo
	// ViewModePomodoro shows the Pomodoro timer screen
	ViewModePomodoro
	// ViewModeDetail shows the todo detail screen
	ViewModeDetail
)

// Model represents the Bubbletea model for the TUI
type Model struct {
	service  *service.TodoService
	todos    []*model.Todo
	cursor   int
	viewMode ViewMode
	input    textinput.Model
	viewport viewport.Model
	message  string
	err      error
	width    int
	height   int
	quitting bool

	// Add todo screen state
	addTodoTitle       string
	addTodoDescription string
	addTodoPriority    model.Priority
	addTodoStep        int // 0: title, 1: description, 2: priority

	// Edit todo screen state
	editTodoID          int64
	editTodoTitle       string
	editTodoDescription string
	editTodoPriority    model.Priority
	editTodoStep        int // 0: title, 1: description, 2: priority

	// Pomodoro timer state
	pomoTodoID      int64 // ID of todo being worked on (0 if general timer)
	pomoSecondsLeft int   // Remaining time in seconds (25 minutes = 1500 seconds)
	pomoRunning     bool  // Whether timer is currently running
	pomoCompleted   bool  // Whether timer has completed and is in alert mode

	// Detail view state
	detailTodoID int64 // ID of todo being displayed in detail view
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
