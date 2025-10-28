package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/syeeel/koto-cli-go/internal/model"
	"github.com/syeeel/koto-cli-go/internal/repository"
)

var (
	// ErrTodoNotFound is returned when a todo is not found
	ErrTodoNotFound = errors.New("todo not found")
	// ErrInvalidTitle is returned when the title is empty or invalid
	ErrInvalidTitle = errors.New("title cannot be empty")
	// ErrInvalidPriority is returned when the priority is invalid
	ErrInvalidPriority = errors.New("invalid priority")
	// ErrInvalidWorkDuration is returned when work duration is invalid (negative or zero)
	ErrInvalidWorkDuration = errors.New("work duration must be positive")
	// ErrFileNotFound is returned when the specified file is not found
	ErrFileNotFound = errors.New("file not found")
	// ErrInvalidJSON is returned when the JSON format is invalid
	ErrInvalidJSON = errors.New("invalid JSON format")
	// ErrExportFailed is returned when export operation fails
	ErrExportFailed = errors.New("export failed")
	// ErrImportFailed is returned when import operation fails
	ErrImportFailed = errors.New("import failed")
)

// TodoService provides business logic for todo operations
type TodoService struct {
	repo repository.TodoRepository
}

// NewTodoService creates a new TodoService
func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

// AddTodo adds a new todo item
func (s *TodoService) AddTodo(ctx context.Context, title, description string, priority model.Priority, dueDate *time.Time) (*model.Todo, error) {
	if err := s.validateTitle(title); err != nil {
		return nil, err
	}

	if err := s.validatePriority(priority); err != nil {
		return nil, err
	}

	now := time.Now()
	todo := &model.Todo{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Status:      model.StatusPending,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// EditTodo edits an existing todo item
func (s *TodoService) EditTodo(ctx context.Context, id int64, title, description string, priority model.Priority, dueDate *time.Time) error {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrTodoNotFound {
			return ErrTodoNotFound
		}
		return err
	}

	if todo == nil {
		return ErrTodoNotFound
	}

	if err := s.validateTitle(title); err != nil {
		return err
	}

	if err := s.validatePriority(priority); err != nil {
		return err
	}

	todo.Title = strings.TrimSpace(title)
	todo.Description = strings.TrimSpace(description)
	todo.Priority = priority
	todo.DueDate = dueDate
	todo.UpdatedAt = time.Now()

	return s.repo.Update(ctx, todo)
}

// DeleteTodo deletes a todo by ID
func (s *TodoService) DeleteTodo(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err == repository.ErrTodoNotFound {
		return ErrTodoNotFound
	}
	return err
}

// CompleteTodo marks a todo as completed
func (s *TodoService) CompleteTodo(ctx context.Context, id int64) error {
	err := s.repo.MarkAsCompleted(ctx, id)
	if err == repository.ErrTodoNotFound {
		return ErrTodoNotFound
	}
	return err
}

// ListTodos returns all todos
func (s *TodoService) ListTodos(ctx context.Context) ([]*model.Todo, error) {
	return s.repo.GetAll(ctx)
}

// ListPendingTodos returns all pending todos
func (s *TodoService) ListPendingTodos(ctx context.Context) ([]*model.Todo, error) {
	return s.repo.GetByStatus(ctx, model.StatusPending)
}

// ListCompletedTodos returns all completed todos
func (s *TodoService) ListCompletedTodos(ctx context.Context) ([]*model.Todo, error) {
	return s.repo.GetByStatus(ctx, model.StatusCompleted)
}

// ExportToJSON exports all todos to a JSON file
func (s *TodoService) ExportToJSON(ctx context.Context, filepath string) error {
	todos, err := s.repo.GetAll(ctx)
	if err != nil {
		return ErrExportFailed
	}

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return ErrExportFailed
	}

	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return ErrExportFailed
	}

	return nil
}

// ImportFromJSON imports todos from a JSON file
func (s *TodoService) ImportFromJSON(ctx context.Context, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return ErrImportFailed
	}

	var todos []*model.Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return ErrInvalidJSON
	}

	// Import each todo (note: this creates new todos, doesn't preserve IDs)
	for _, todo := range todos {
		// Reset ID to create as new todo
		todo.ID = 0
		todo.CreatedAt = time.Now()
		todo.UpdatedAt = time.Now()

		if err := s.repo.Create(ctx, todo); err != nil {
			return ErrImportFailed
		}
	}

	return nil
}

// validateTitle validates the todo title
func (s *TodoService) validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrInvalidTitle
	}
	return nil
}

// validatePriority validates the priority value
func (s *TodoService) validatePriority(priority model.Priority) error {
	if priority < model.PriorityLow || priority > model.PriorityHigh {
		return ErrInvalidPriority
	}
	return nil
}

// AddWorkDuration adds work duration (in minutes) to a todo
// Returns an error if the todo doesn't exist or if the duration is invalid
func (s *TodoService) AddWorkDuration(ctx context.Context, id int64, minutes int) error {
	// Validate work duration (must be positive)
	if minutes <= 0 {
		return ErrInvalidWorkDuration
	}

	// Check if todo exists
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrTodoNotFound {
			return ErrTodoNotFound
		}
		return err
	}

	if todo == nil {
		return ErrTodoNotFound
	}

	// Add work duration
	return s.repo.AddWorkDuration(ctx, id, minutes)
}
