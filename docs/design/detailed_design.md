# koto - ToDo Management CLI Tool Detailed Design Document

## 1. Directory Structure

```
koto/
├── cmd/
│   └── koto/
│       └── main.go              # Entry point
├── internal/
│   ├── model/
│   │   └── todo.go              # Todo data model
│   ├── repository/
│   │   ├── repository.go        # Repository interface
│   │   └── sqlite.go            # SQLite implementation
│   ├── service/
│   │   └── todo_service.go      # Business logic
│   ├── tui/
│   │   ├── model.go             # Bubbletea Model
│   │   ├── commands.go          # Command parser
│   │   ├── views.go             # View rendering
│   │   ├── styles.go            # Lipgloss styles
│   │   └── update.go            # Update function
│   └── config/
│       └── config.go            # Configuration management
├── pkg/
│   └── utils/
│       ├── time.go              # Time-related utilities
│       └── validation.go        # Validation
├── migrations/
│   └── 001_init.sql             # DB migration
├── docs/
│   ├── basic_design.md          # Basic design document
│   └── detailed_design.md       # Detailed design document (this document)
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── Makefile
```

## 2. Package Structure

### 2.1 cmd/koto
- Application entry point
- Configuration initialization
- Launch TUI application

### 2.2 internal/model
- Data model definition
- Plain Old Go Object (POGO) with no domain logic

### 2.3 internal/repository
- Data persistence abstraction layer
- SQLite implementation
- Makes it easy to support other storage in the future

### 2.4 internal/service
- Business logic layer
- Data manipulation using Repository
- Validation execution

### 2.5 internal/tui
- TUI implementation using Bubbletea
- Model-View-Update pattern
- Command parser
- Styling

### 2.6 internal/config
- Application configuration
- Database path management
- Configuration file reading (optional)

### 2.7 pkg/utils
- General utility functions
- Functions that can be used from outside the project

## 3. Data Model Details

### 3.1 Todo Structure

```go
package model

import "time"

type TodoStatus int

const (
    StatusPending TodoStatus = iota
    StatusCompleted
)

type Priority int

const (
    PriorityLow Priority = iota
    PriorityMedium
    PriorityHigh
)

type Todo struct {
    ID           int64      `db:"id"`
    Title        string     `db:"title"`
    Description  string     `db:"description"`
    Status       TodoStatus `db:"status"`
    Priority     Priority   `db:"priority"`
    DueDate      *time.Time `db:"due_date"`
    WorkDuration int        `db:"work_duration"` // Cumulative work time (minutes)
    CreatedAt    time.Time  `db:"created_at"`
    UpdatedAt    time.Time  `db:"updated_at"`
}

func (t Todo) IsCompleted() bool {
    return t.Status == StatusCompleted
}

func (t Todo) IsPending() bool {
    return t.Status == StatusPending
}

func (t Todo) IsOverdue() bool {
    if t.DueDate == nil {
        return false
    }
    return time.Now().After(*t.DueDate) && t.IsPending()
}
```

### 3.2 Database Schema

```sql
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    priority INTEGER NOT NULL DEFAULT 0,
    due_date DATETIME,
    work_duration INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_due_date ON todos(due_date);
CREATE INDEX idx_todos_created_at ON todos(created_at);
```

## 4. Repository Layer Details

### 4.1 Repository Interface

```go
package repository

import (
    "context"
    "github.com/yourusername/koto/internal/model"
)

type TodoRepository interface {
    // Create creates a new todo item
    Create(ctx context.Context, todo *model.Todo) error

    // GetByID retrieves a todo by ID
    GetByID(ctx context.Context, id int64) (*model.Todo, error)

    // GetAll retrieves all todos
    GetAll(ctx context.Context) ([]*model.Todo, error)

    // GetByStatus retrieves todos by status
    GetByStatus(ctx context.Context, status model.TodoStatus) ([]*model.Todo, error)

    // Update updates a todo
    Update(ctx context.Context, todo *model.Todo) error

    // Delete deletes a todo by ID
    Delete(ctx context.Context, id int64) error

    // MarkAsCompleted marks a todo as completed
    MarkAsCompleted(ctx context.Context, id int64) error

    // AddWorkDuration adds work duration (in minutes) to a todo
    AddWorkDuration(ctx context.Context, id int64, minutes int) error

    // Close closes the repository connection
    Close() error
}
```

### 4.2 SQLite Implementation

Using modernc.org/sqlite enables a Pure Go implementation without CGO.

```go
package repository

import (
    "context"
    "database/sql"
    "time"

    _ "modernc.org/sqlite"  // Pure Go SQLite driver
    "github.com/yourusername/koto/internal/model"
)

type SQLiteRepository struct {
    db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
    // Using modernc.org/sqlite (Pure Go, no CGO required)
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    // Set file permissions (security)
    if err := os.Chmod(dbPath, 0600); err != nil {
        return nil, err
    }

    // Initialize schema
    if err := initSchema(db); err != nil {
        return nil, err
    }

    return &SQLiteRepository{db: db}, nil
}

// Implementation of each method...
```

**Notes:**
- Can be built with `CGO_ENABLED=0`
- Easy cross-compilation
- Simple dependencies (Pure Go)

## 5. Service Layer Details

### 5.1 TodoService

```go
package service

import (
    "context"
    "errors"
    "time"

    "github.com/yourusername/koto/internal/model"
    "github.com/yourusername/koto/internal/repository"
)

var (
    ErrTodoNotFound    = errors.New("todo not found")
    ErrInvalidTitle    = errors.New("title cannot be empty")
    ErrInvalidPriority = errors.New("invalid priority")
)

type TodoService struct {
    repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
    return &TodoService{repo: repo}
}

func (s *TodoService) AddTodo(ctx context.Context, title, description string, priority model.Priority, dueDate *time.Time) (*model.Todo, error) {
    if err := s.validateTitle(title); err != nil {
        return nil, err
    }

    todo := &model.Todo{
        Title:       title,
        Description: description,
        Status:      model.StatusPending,
        Priority:    priority,
        DueDate:     dueDate,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := s.repo.Create(ctx, todo); err != nil {
        return nil, err
    }

    return todo, nil
}

func (s *TodoService) EditTodo(ctx context.Context, id int64, title, description string, priority model.Priority, dueDate *time.Time) error {
    todo, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    if todo == nil {
        return ErrTodoNotFound
    }

    if err := s.validateTitle(title); err != nil {
        return err
    }

    todo.Title = title
    todo.Description = description
    todo.Priority = priority
    todo.DueDate = dueDate
    todo.UpdatedAt = time.Now()

    return s.repo.Update(ctx, todo)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int64) error {
    return s.repo.Delete(ctx, id)
}

func (s *TodoService) CompleteTodo(ctx context.Context, id int64) error {
    return s.repo.MarkAsCompleted(ctx, id)
}

func (s *TodoService) ListTodos(ctx context.Context) ([]*model.Todo, error) {
    return s.repo.GetAll(ctx)
}

func (s *TodoService) ListPendingTodos(ctx context.Context) ([]*model.Todo, error) {
    return s.repo.GetByStatus(ctx, model.StatusPending)
}

func (s *TodoService) ListCompletedTodos(ctx context.Context) ([]*model.Todo, error) {
    return s.repo.GetByStatus(ctx, model.StatusCompleted)
}

func (s *TodoService) validateTitle(title string) error {
    if title == "" {
        return ErrInvalidTitle
    }
    return nil
}

func (s *TodoService) ExportToJSON(ctx context.Context, filepath string) error {
    todos, err := s.repo.GetAll(ctx)
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(todos, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filepath, data, 0600)
}

func (s *TodoService) ImportFromJSON(ctx context.Context, filepath string) error {
    data, err := os.ReadFile(filepath)
    if err != nil {
        return err
    }

    var todos []*model.Todo
    if err := json.Unmarshal(data, &todos); err != nil {
        return err
    }

    // Consider implementing user confirmation before deleting existing data
    for _, todo := range todos {
        if err := s.repo.Create(ctx, todo); err != nil {
            return err
        }
    }

    return nil
}

func (s *TodoService) AddWorkDuration(ctx context.Context, id int64, minutes int) error {
    // Check if the ToDo exists
    todo, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if todo == nil {
        return ErrTodoNotFound
    }

    // Add work duration
    return s.repo.AddWorkDuration(ctx, id, minutes)
}
```

## 6. TUI Layer Detailed Design

### 6.1 Bubbletea Model

```go
package tui

import (
    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourusername/koto/internal/model"
    "github.com/yourusername/koto/internal/service"
)

type ViewMode int

const (
    ViewModeList ViewMode = iota
    ViewModeAdd
    ViewModeEdit
    ViewModeDelete
    ViewModeHelp
    ViewModePomodoro
    ViewModeDetail  // Task detail display screen
)

type Model struct {
    service         *service.TodoService
    todos           []*model.Todo
    cursor          int
    viewMode        ViewMode
    input           textinput.Model
    message         string
    err             error
    width           int
    height          int
    selectedID      int64
    pomodoroTodoID  *int64        // ToDo ID linked to Pomodoro timer (nil if not linked)
    pomodoroStarted time.Time     // Pomodoro start time
    pomodoroDuration time.Duration // Pomodoro duration (default 25 minutes)
    detailTodoID    int64         // ToDo ID being displayed in detail
}

func NewModel(service *service.TodoService) Model {
    ti := textinput.New()
    ti.Placeholder = "Enter a command (/help for help)"
    ti.Focus()
    ti.CharLimit = 500
    ti.Width = 80

    return Model{
        service:  service,
        todos:    []*model.Todo{},
        cursor:   0,
        viewMode: ViewModeList,
        input:    ti,
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(textinput.Blink, loadTodos(m.service))
}
```

### 6.2 Update Function

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "esc":
            if m.viewMode != ViewModeList {
                m.viewMode = ViewModeList
                m.input.SetValue("")
                return m, nil
            }
            return m, tea.Quit

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
        }

    case todosLoadedMsg:
        m.todos = msg.todos
        m.err = msg.err
        return m, nil
    }

    m.input, cmd = m.input.Update(msg)
    return m, cmd
}

func (m *Model) handleEnter() (tea.Model, tea.Cmd) {
    value := m.input.Value()
    m.input.SetValue("")

    return m, parseAndExecuteCommand(m.service, value)
}
```

### 6.3 View Function

```go
func (m Model) View() string {
    switch m.viewMode {
    case ViewModeList:
        return m.renderListView()
    case ViewModeAdd:
        return m.renderAddView()
    case ViewModeEdit:
        return m.renderEditView()
    case ViewModeDelete:
        return m.renderDeleteView()
    case ViewModeHelp:
        return m.renderHelpView()
    case ViewModePomodoro:
        return m.renderPomodoroView()
    case ViewModeDetail:
        return m.renderDetailView()
    default:
        return m.renderListView()
    }
}

func (m Model) renderListView() string {
    var s string

    s += titleStyle.Render("📝 koto - ToDo Manager") + "\n\n"

    if len(m.todos) == 0 {
        s += emptyStyle.Render("No ToDos. Add a new ToDo with /add.") + "\n"
    } else {
        // Header row
        s += headerStyle.Render("No.  Title                Description          Created At") + "\n"
        s += separatorStyle.Render(strings.Repeat("─", 80)) + "\n"

        for i, todo := range m.todos {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }

            // No. (ID)
            no := fmt.Sprintf("%-4d", todo.ID)

            // Title (max 20 chars, truncate)
            title := truncateString(todo.Title, 20)

            // Description (max 20 chars, truncate)
            desc := truncateString(todo.Description, 20)

            // Created At (date only)
            created := todo.CreatedAt.Format("2006-01-02")

            // Build the line
            line := fmt.Sprintf("%s %s %s %s %s", cursor, no, title, desc, created)

            if m.cursor == i {
                s += selectedStyle.Render(line) + "\n"
            } else {
                s += line + "\n"
            }

            // Separator line
            s += separatorStyle.Render(strings.Repeat("─", 80)) + "\n"
        }
    }

    s += "\n" + m.input.View() + "\n"

    if m.message != "" {
        s += "\n" + messageStyle.Render(m.message) + "\n"
    }

    if m.err != nil {
        s += "\n" + errorStyle.Render("Error: "+m.err.Error()) + "\n"
    }

    s += "\n" + helpStyle.Render("Usage: /help | Exit: Ctrl+C")

    return s
}

func (m Model) renderPomodoroView() string {
    var s string

    // Calculate remaining time
    elapsed := time.Since(m.pomodoroStarted)
    remaining := m.pomodoroDuration - elapsed

    if remaining < 0 {
        remaining = 0
    }

    minutes := int(remaining.Minutes())
    seconds := int(remaining.Seconds()) % 60

    s += pomodoroTitleStyle.Render("🍅 Pomodoro Timer") + "\n\n"

    // Timer display (large)
    timerText := fmt.Sprintf("%02d:%02d", minutes, seconds)
    s += pomodoroTimerStyle.Render(timerText) + "\n\n"

    // If linked to a ToDo, display the title
    if m.pomodoroTodoID != nil {
        for _, todo := range m.todos {
            if todo.ID == *m.pomodoroTodoID {
                s += pomodoroTaskStyle.Render(fmt.Sprintf("Working on: %s", todo.Title)) + "\n\n"
                break
            }
        }
    } else {
        s += pomodoroTaskStyle.Render("Free timer mode") + "\n\n"
    }

    // When timer ends
    if remaining == 0 {
        s += pomodoroCompleteStyle.Render("🎉 Pomodoro completed!") + "\n"
        s += "Press Enter to return to main screen\n"
    } else {
        s += helpStyle.Render("Press Esc to cancel")
    }

    return s
}

func (m Model) renderDetailView() string {
    var s string

    // Get the ToDo to display in detail
    var targetTodo *model.Todo
    for _, todo := range m.todos {
        if todo.ID == m.detailTodoID {
            targetTodo = todo
            break
        }
    }

    if targetTodo == nil {
        s += errorStyle.Render("Specified ToDo not found") + "\n"
        s += helpStyle.Render("Press Esc to return to main screen")
        return s
    }

    // Title
    s += titleStyle.Render(fmt.Sprintf("📋 ToDo Details #%d", targetTodo.ID)) + "\n\n"

    // Title display
    s += headerStyle.Render(" Title ") + "\n"
    s += todoDetailFieldStyle.Render(targetTodo.Title) + "\n\n"

    // Description display
    s += headerStyle.Render(" Description ") + "\n"
    if targetTodo.Description != "" {
        s += todoDetailFieldStyle.Render(targetTodo.Description) + "\n\n"
    } else {
        s += emptyStyle.Render("(No description)") + "\n\n"
    }

    // Status display
    s += headerStyle.Render(" Status ") + "\n"
    status := "Incomplete"
    if targetTodo.IsCompleted() {
        status = "Completed"
    }
    s += todoDetailFieldStyle.Render(status) + "\n\n"

    // Priority display
    s += headerStyle.Render(" Priority ") + "\n"
    priorityStr := ""
    switch targetTodo.Priority {
    case model.PriorityHigh:
        priorityStr = "High"
    case model.PriorityMedium:
        priorityStr = "Medium"
    case model.PriorityLow:
        priorityStr = "Low"
    }
    s += todoDetailFieldStyle.Render(priorityStr) + "\n\n"

    // Cumulative work time display
    s += headerStyle.Render(" Cumulative Work Time ") + "\n"
    if targetTodo.WorkDuration > 0 {
        hours := targetTodo.WorkDuration / 60
        minutes := targetTodo.WorkDuration % 60
        workDurationStr := fmt.Sprintf("%dh %dm", hours, minutes)
        s += todoDetailFieldStyle.Render(workDurationStr) + "\n\n"
    } else {
        s += emptyStyle.Render("(No record)") + "\n\n"
    }

    // Due date display
    s += headerStyle.Render(" Due Date ") + "\n"
    if targetTodo.DueDate != nil {
        dueDateStr := targetTodo.DueDate.Format("2006-01-02 15:04")
        s += todoDetailFieldStyle.Render(dueDateStr)
        if targetTodo.IsOverdue() {
            s += " " + errorStyle.Render("(Overdue)")
        }
        s += "\n\n"
    } else {
        s += emptyStyle.Render("(No due date)") + "\n\n"
    }

    // Created/Updated times
    s += headerStyle.Render(" Created At ") + "\n"
    s += todoDetailFieldStyle.Render(targetTodo.CreatedAt.Format("2006-01-02 15:04:05")) + "\n\n"

    s += headerStyle.Render(" Updated At ") + "\n"
    s += todoDetailFieldStyle.Render(targetTodo.UpdatedAt.Format("2006-01-02 15:04:05")) + "\n\n"

    // Help text
    s += "\n" + helpStyle.Render("Press Esc to return to main screen | e to edit | d to complete")

    return s
}
```

### 6.4 Command Parser

```go
package tui

import (
    "context"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourusername/koto/internal/service"
)

type commandExecutedMsg struct {
    message string
    err     error
}

type todosLoadedMsg struct {
    todos []*model.Todo
    err   error
}

func parseAndExecuteCommand(svc *service.TodoService, input string) tea.Cmd {
    return func() tea.Msg {
        input = strings.TrimSpace(input)

        if !strings.HasPrefix(input, "/") {
            return commandExecutedMsg{
                err: errors.New("command must start with /"),
            }
        }

        parts := strings.Fields(input)
        command := parts[0]
        args := parts[1:]

        ctx := context.Background()

        switch command {
        case "/add":
            return handleAddCommand(ctx, svc, args)
        case "/edit":
            return handleEditCommand(ctx, svc, args)
        case "/delete":
            return handleDeleteCommand(ctx, svc, args)
        case "/done":
            return handleDoneCommand(ctx, svc, args)
        case "/list":
            return handleListCommand(ctx, svc)
        case "/export":
            return handleExportCommand(ctx, svc, args)
        case "/import":
            return handleImportCommand(ctx, svc, args)
        case "/pomo":
            return handlePomodoroCommand(ctx, svc, args)
        case "/help":
            return commandExecutedMsg{message: getHelpText()}
        case "/quit":
            return tea.Quit()
        default:
            return commandExecutedMsg{
                err: errors.New("unknown command: " + command),
            }
        }
    }
}

func loadTodos(svc *service.TodoService) tea.Cmd {
    return func() tea.Msg {
        todos, err := svc.ListTodos(context.Background())
        return todosLoadedMsg{todos: todos, err: err}
    }
}

// Implementation of each command handler...
```

### 6.5 Style Definitions

```go
package tui

import "github.com/charmbracelet/lipgloss"

var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        MarginBottom(1)

    emptyStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Italic(true)

    messageStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Bold(true)

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true)

    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        MarginTop(1)

    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("212")).
        Bold(true)

    pomodoroTitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("196")).
        MarginBottom(1)

    pomodoroTimerStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("226")).
        FontSize(48).  // Large font (pseudo-representation in TUI)
        Align(lipgloss.Center)

    pomodoroTaskStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Italic(true)

    pomodoroCompleteStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("46"))
)
```

## 7. Command Specification Details

### 7.1 /add - Add ToDo

**Syntax:**
```
/add <title> [--desc=<description>] [--priority=<low|medium|high>] [--due=<YYYY-MM-DD>]
```

**Examples:**
```
/add Write report
/add Write report --desc=Monthly report --priority=high --due=2024-12-31
```

**Processing Flow:**
1. Parse command
2. Validation (title required)
3. Call Service layer's AddTodo
4. Display success message
5. Reload list

### 7.2 /edit - Edit ToDo

**Syntax:**
```
/edit <ID> [--title=<new title>] [--desc=<new description>] [--priority=<low|medium|high>] [--due=<YYYY-MM-DD>]
```

**Examples:**
```
/edit 1 --title=New title
/edit 2 --priority=high --due=2024-12-25
```

**Processing Flow:**
1. Search for Todo by ID
2. Check existence
3. Update specified fields
4. Call Service layer's EditTodo
5. Display success message
6. Reload list

### 7.3 /delete - Delete ToDo

**Syntax:**
```
/delete <ID>
```

**Examples:**
```
/delete 1
```

**Processing Flow:**
1. Search for Todo by ID
2. Check existence
3. Display deletion confirmation prompt (optional)
4. Call Service layer's DeleteTodo
5. Display success message
6. Reload list

### 7.4 /done - Complete ToDo

**Syntax:**
```
/done <ID>
```

**Examples:**
```
/done 1
```

**Processing Flow:**
1. Search for Todo by ID
2. Check existence
3. Call Service layer's CompleteTodo
4. Display success message
5. Reload list

### 7.5 /list - Display List

**Syntax:**
```
/list [--status=<pending|completed|all>]
```

**Examples:**
```
/list
/list --status=pending
/list --status=completed
```

**Processing Flow:**
1. Parse status filter
2. Call corresponding Service layer method
3. Display list

### 7.6 /help - Display Help

**Syntax:**
```
/help
```

**Processing Flow:**
1. Display help text

### 7.7 /export - JSON Export

**Syntax:**
```
/export [file path]
```

**Examples:**
```
/export
/export ~/backups/todos-2024-12-01.json
```

**Processing Flow:**
1. Parse export destination file path (default: `~/.koto/export.json`)
2. Call Service layer's ExportToJSON
3. Write all ToDo data to file in JSON format
4. Display success message

**Export Format:**
```json
[
  {
    "id": 1,
    "title": "Write report",
    "description": "Monthly report",
    "status": 0,
    "priority": 2,
    "due_date": "2024-12-31T23:59:59Z",
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
  }
]
```

### 7.8 /import - JSON Import

**Syntax:**
```
/import <file path>
```

**Examples:**
```
/import ~/backups/todos-2024-12-01.json
```

**Processing Flow:**
1. Parse import source file path
2. Check file existence
3. Validate JSON format
4. Display confirmation prompt (regarding duplicates with existing data)
5. Call Service layer's ImportFromJSON
6. Display success message (including number of imported items)
7. Reload list

**Notes:**
- Confirm with user how to handle duplicate IDs (overwrite or skip)
- Display progress for large data imports

### 7.9 /pomo - Pomodoro Timer

**Syntax:**
```
/pomo [ToDo ID]
```

**Examples:**
```
/pomo              # Free timer mode (25 minutes)
/pomo 1            # Start 25-minute timer linked to ToDo ID 1
```

**Processing Flow:**
1. Check arguments (if ToDo ID is specified, verify existence)
2. Transition to Pomodoro mode
3. Start timer (25 minutes = 1500 seconds)
4. Update screen every second (using tea.Tick)
5. When timer ends
   - Display alarm
   - If ToDo ID is specified, record work time (add 25 minutes)
   - Wait for Enter key
6. Return to main screen on Enter key

**Cancellation:**
- Press Esc to cancel timer and return to main screen
- No work time recorded on cancellation

**Notes:**
- Dedicated screen is displayed during timer execution
- Background execution not supported
- Work time only recorded when timer completes

### 7.10 /quit - Exit App

**Syntax:**
```
/quit
```

**Processing Flow:**
1. Close database connection
2. Exit application

## 8. Error Handling

### 8.1 Error Types

| Error Type | Description | Response |
|-----------|-------------|----------|
| Validation Error | Invalid input value | Display error message and prompt re-entry |
| Database Error | DB operation failed | Output error log and notify user |
| Non-existent ID | Specified ID Todo not found | Display error message |
| Command Parse Error | Invalid command | Display help and guide to correct syntax |
| File I/O Error | Export/import failed | Guide to check file path or permissions |
| JSON Parse Error | Invalid JSON format | Guide to specify correctly exported JSON file |

### 8.2 Error Message Examples

```go
var (
    ErrMessages = map[error]string{
        service.ErrTodoNotFound:    "Specified ToDo not found",
        service.ErrInvalidTitle:    "Title is required",
        service.ErrInvalidPriority: "Priority must be low, medium, or high",
        service.ErrFileNotFound:    "File not found",
        service.ErrInvalidJSON:     "Invalid JSON format",
        service.ErrExportFailed:    "Export failed",
        service.ErrImportFailed:    "Import failed",
    }
)
```

## 9. Data Persistence

### 9.1 Database File Location

- **Path**: `~/.koto/koto.db`
- **Permissions**: 0600 (read/write only by owner)

### 9.2 Initialization Process

```go
func initDatabase() (*sql.DB, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    kotoDir := filepath.Join(homeDir, ".koto")
    if err := os.MkdirAll(kotoDir, 0700); err != nil {
        return nil, err
    }

    dbPath := filepath.Join(kotoDir, "koto.db")
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    return db, nil
}
```

## 10. Testing Strategy

### 10.1 Unit Tests

- **Target**: Functions and methods of each package
- **Tools**: Go standard testing package
- **Coverage Goal**: 80% or higher

**Example:**
```go
func TestTodoService_AddTodo(t *testing.T) {
    // Setup
    repo := repository.NewMockRepository()
    svc := service.NewTodoService(repo)

    // Execute
    todo, err := svc.AddTodo(context.Background(), "Test Todo", "", model.PriorityMedium, nil)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, todo)
    assert.Equal(t, "Test Todo", todo.Title)
}
```

### 10.2 Integration Tests

- **Target**: Integration of Repository layer and database
- **Tools**: Go testing + In-memory SQLite

### 10.3 E2E Tests

- **Target**: Application flow as a whole
- **Tools**: Go testing + Bubbletea test utilities

## 11. Build and Release

### 11.1 Makefile

```makefile
.PHONY: build test clean install release-local

# Pure Go build (no CGO required)
build:
	CGO_ENABLED=0 go build -o bin/koto ./cmd/koto

# Cross-compilation build
build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/koto-linux-amd64 ./cmd/koto
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/koto-darwin-amd64 ./cmd/koto
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/koto-darwin-arm64 ./cmd/koto
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/koto-windows-amd64.exe ./cmd/koto

test:
	go test -v -cover ./...

clean:
	rm -rf bin/

install: build
	cp bin/koto $(GOPATH)/bin/koto

lint:
	golangci-lint run

release:
	goreleaser release --clean
```

### 11.2 GoReleaser Configuration

```yaml
# .goreleaser.yml
project_name: koto

builds:
  - id: koto
    main: ./cmd/koto
    binary: koto
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
```

## 12. Achieving Data Portability

### 12.1 Export Feature
While SQLite is in binary format, the `/export` command achieves the following:
- Conversion to human-readable JSON format
- Easy migration to other environments
- Backup creation
- Integration with other tools

### 12.2 Import Feature
The `/import` command achieves the following:
- Restore from backup
- Data migration from other environments
- Import ToDo data created by external tools

### 12.3 Data Integration Examples
```bash
# Environment A: Export data
/export ~/todos-backup.json

# Environment B: Import data
/import ~/todos-backup.json

# Processing with scripts is also possible
cat ~/todos-backup.json | jq '.[] | select(.priority == 2)' > high-priority.json
/import high-priority.json
```

## 13. Future Extensions

### 13.1 Phase 2 (v2.0)
- Tag feature (assign multiple tags)
- Advanced filtering and search (full-text search, multiple conditions)
- Category classification

### 13.2 Phase 3 (v3.0)
- Recurring tasks
- CSV export
- Configuration file support (YAML/TOML)
- Subtask feature

### 13.3 Phase 4 (v4.0 and beyond)
- Cloud sync
- Team sharing
- Web UI
- Mobile app integration
