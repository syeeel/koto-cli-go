package tui

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/syeeel/koto-cli-go/internal/model"
	"github.com/syeeel/koto-cli-go/internal/service"
)

// Message types for Bubbletea

// commandExecutedMsg is sent when a command has been executed
type commandExecutedMsg struct {
	message string
	err     error
}

// todosLoadedMsg is sent when todos have been loaded
type todosLoadedMsg struct {
	todos []*model.Todo
	err   error
}

// pomodoroTickMsg is sent every second when the timer is running
type pomodoroTickMsg struct{}

// pomodoroCompleteMsg is sent when the timer reaches zero
type pomodoroCompleteMsg struct {
	todoID int64 // ID of todo that was worked on (0 if general timer)
}

// parseAndExecuteCommand parses and executes a command
func parseAndExecuteCommand(svc *service.TodoService, input string) tea.Cmd {
	return func() tea.Msg {
		input = strings.TrimSpace(input)

		// Check if input starts with /
		if !strings.HasPrefix(input, "/") {
			return commandExecutedMsg{
				err: errors.New("commands must start with /"),
			}
		}

		// Parse command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return commandExecutedMsg{
				err: errors.New("no command provided"),
			}
		}

		command := parts[0]
		args := parts[1:]

		ctx := context.Background()

		// Execute command
		switch command {
		case "/delete":
			return handleDeleteCommand(ctx, svc, args)
		case "/done":
			return handleDoneCommand(ctx, svc, args)
		case "/list":
			return handleListCommand(ctx, svc, args)
		case "/export":
			return handleExportCommand(ctx, svc, args)
		case "/import":
			return handleImportCommand(ctx, svc, args)
		case "/help":
			return commandExecutedMsg{message: "Press '?' to view help"}
		case "/exit":
			return tea.Quit()
		default:
			return commandExecutedMsg{
				err: fmt.Errorf("unknown command: %s", command),
			}
		}
	}
}

// loadTodos loads all todos from the service
func loadTodos(svc *service.TodoService) tea.Cmd {
	return func() tea.Msg {
		todos, err := svc.ListTodos(context.Background())
		return todosLoadedMsg{todos: todos, err: err}
	}
}

// handleDeleteCommand handles the /delete command
func handleDeleteCommand(ctx context.Context, svc *service.TodoService, args []string) commandExecutedMsg {
	if len(args) != 1 {
		return commandExecutedMsg{err: errors.New("usage: /delete <id>")}
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return commandExecutedMsg{err: errors.New("invalid todo ID")}
	}

	err = svc.DeleteTodo(ctx, id)
	if err != nil {
		return commandExecutedMsg{err: err}
	}

	return commandExecutedMsg{message: fmt.Sprintf("Deleted todo #%d", id)}
}

// handleDoneCommand handles the /done command
func handleDoneCommand(ctx context.Context, svc *service.TodoService, args []string) commandExecutedMsg {
	if len(args) != 1 {
		return commandExecutedMsg{err: errors.New("usage: /done <id>")}
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return commandExecutedMsg{err: errors.New("invalid todo ID")}
	}

	err = svc.CompleteTodo(ctx, id)
	if err != nil {
		return commandExecutedMsg{err: err}
	}

	return commandExecutedMsg{message: fmt.Sprintf("Marked todo #%d as completed", id)}
}

// handleListCommand handles the /list command
func handleListCommand(ctx context.Context, svc *service.TodoService, args []string) commandExecutedMsg {
	// Parse status filter
	var status *model.TodoStatus
	for _, arg := range args {
		if strings.HasPrefix(arg, "--status=") {
			statusStr := strings.TrimPrefix(arg, "--status=")
			switch strings.ToLower(statusStr) {
			case "pending":
				s := model.StatusPending
				status = &s
			case "completed":
				s := model.StatusCompleted
				status = &s
			case "all":
				status = nil
			default:
				return commandExecutedMsg{err: errors.New("invalid status (use: pending, completed, all)")}
			}
		}
	}

	// This command just triggers a reload, the actual display is handled by the view
	var todos []*model.Todo
	var err error

	if status == nil {
		todos, err = svc.ListTodos(ctx)
	} else if *status == model.StatusPending {
		todos, err = svc.ListPendingTodos(ctx)
	} else {
		todos, err = svc.ListCompletedTodos(ctx)
	}

	if err != nil {
		return commandExecutedMsg{err: err}
	}

	return commandExecutedMsg{message: fmt.Sprintf("Showing %d todos", len(todos))}
}

// handleExportCommand handles the /export command
func handleExportCommand(ctx context.Context, svc *service.TodoService, args []string) commandExecutedMsg {
	filepath := "todos_export.json"
	if len(args) > 0 {
		filepath = args[0]
	}

	err := svc.ExportToJSON(ctx, filepath)
	if err != nil {
		return commandExecutedMsg{err: err}
	}

	return commandExecutedMsg{message: fmt.Sprintf("Exported todos to %s", filepath)}
}

// handleImportCommand handles the /import command
func handleImportCommand(ctx context.Context, svc *service.TodoService, args []string) commandExecutedMsg {
	if len(args) != 1 {
		return commandExecutedMsg{err: errors.New("usage: /import <filepath>")}
	}

	filepath := args[0]

	err := svc.ImportFromJSON(ctx, filepath)
	if err != nil {
		return commandExecutedMsg{err: err}
	}

	return commandExecutedMsg{message: fmt.Sprintf("Imported todos from %s", filepath)}
}

// tickPomodoro creates a command that waits 1 second and sends a tick message
func tickPomodoro() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return pomodoroTickMsg{}
	})
}

// completePomodoroWithRecording handles timer completion and records work duration
func completePomodoroWithRecording(svc *service.TodoService, todoID int64) tea.Cmd {
	return func() tea.Msg {
		// If this was a task-specific timer, record the work duration
		if todoID > 0 {
			ctx := context.Background()
			err := svc.AddWorkDuration(ctx, todoID, 25) // 25 minutes
			if err != nil {
				return commandExecutedMsg{
					err: fmt.Errorf("failed to record work duration: %w", err),
				}
			}
		}
		return pomodoroCompleteMsg{todoID: todoID}
	}
}

// recordPartialPomodoro records partial work duration when timer is stopped early
func recordPartialPomodoro(svc *service.TodoService, todoID int64, minutes int) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := svc.AddWorkDuration(ctx, todoID, minutes)
		if err != nil {
			return commandExecutedMsg{
				err: fmt.Errorf("failed to record work duration: %w", err),
			}
		}
		return commandExecutedMsg{
			message: fmt.Sprintf("Pomodoro stopped. %d minutes recorded for todo #%d", minutes, todoID),
		}
	}
}
