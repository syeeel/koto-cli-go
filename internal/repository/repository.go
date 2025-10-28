package repository

import (
	"context"

	"github.com/syeeel/koto-cli-go/internal/model"
)

// TodoRepository defines the interface for todo data persistence
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
