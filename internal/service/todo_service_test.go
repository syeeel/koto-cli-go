package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syeeel/koto-cli-go/internal/model"
	"github.com/syeeel/koto-cli-go/internal/repository"
)

// mockRepository is a simple in-memory implementation for testing
type mockRepository struct {
	todos  map[int64]*model.Todo
	nextID int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		todos:  make(map[int64]*model.Todo),
		nextID: 1,
	}
}

func (m *mockRepository) Create(ctx context.Context, todo *model.Todo) error {
	todo.ID = m.nextID
	m.nextID++
	m.todos[todo.ID] = todo
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	todo, exists := m.todos[id]
	if !exists {
		return nil, repository.ErrTodoNotFound
	}
	return todo, nil
}

func (m *mockRepository) GetAll(ctx context.Context) ([]*model.Todo, error) {
	todos := make([]*model.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (m *mockRepository) GetByStatus(ctx context.Context, status model.TodoStatus) ([]*model.Todo, error) {
	todos := make([]*model.Todo, 0)
	for _, todo := range m.todos {
		if todo.Status == status {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (m *mockRepository) Update(ctx context.Context, todo *model.Todo) error {
	if _, exists := m.todos[todo.ID]; !exists {
		return repository.ErrTodoNotFound
	}
	m.todos[todo.ID] = todo
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.todos[id]; !exists {
		return repository.ErrTodoNotFound
	}
	delete(m.todos, id)
	return nil
}

func (m *mockRepository) MarkAsCompleted(ctx context.Context, id int64) error {
	todo, exists := m.todos[id]
	if !exists {
		return repository.ErrTodoNotFound
	}
	todo.Status = model.StatusCompleted
	todo.UpdatedAt = time.Now()
	return nil
}

func (m *mockRepository) AddWorkDuration(ctx context.Context, id int64, minutes int) error {
	todo, exists := m.todos[id]
	if !exists {
		return repository.ErrTodoNotFound
	}
	todo.WorkDuration += minutes
	todo.UpdatedAt = time.Now()
	return nil
}

func (m *mockRepository) Close() error {
	return nil
}

// Tests

func TestNewTodoService(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)

	if svc == nil {
		t.Fatal("expected service to be created, got nil")
	}

	if svc.repo == nil {
		t.Fatal("expected repository to be set, got nil")
	}
}

func TestTodoService_AddTodo(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	tests := []struct {
		name        string
		title       string
		description string
		priority    model.Priority
		dueDate     *time.Time
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid todo",
			title:       "Test Todo",
			description: "Test Description",
			priority:    model.PriorityMedium,
			dueDate:     nil,
			wantErr:     false,
		},
		{
			name:        "empty title",
			title:       "",
			description: "Test Description",
			priority:    model.PriorityMedium,
			dueDate:     nil,
			wantErr:     true,
			expectedErr: ErrInvalidTitle,
		},
		{
			name:        "whitespace only title",
			title:       "   ",
			description: "Test Description",
			priority:    model.PriorityMedium,
			dueDate:     nil,
			wantErr:     true,
			expectedErr: ErrInvalidTitle,
		},
		{
			name:        "invalid priority too low",
			title:       "Test Todo",
			description: "Test Description",
			priority:    model.Priority(-1),
			dueDate:     nil,
			wantErr:     true,
			expectedErr: ErrInvalidPriority,
		},
		{
			name:        "invalid priority too high",
			title:       "Test Todo",
			description: "Test Description",
			priority:    model.Priority(10),
			dueDate:     nil,
			wantErr:     true,
			expectedErr: ErrInvalidPriority,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := svc.AddTodo(ctx, tt.title, tt.description, tt.priority, tt.dueDate)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if todo == nil {
					t.Error("expected todo to be created, got nil")
				}
				if todo.ID == 0 {
					t.Error("expected todo ID to be set")
				}
			}
		})
	}
}

func TestTodoService_EditTodo(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create a todo first
	todo, _ := svc.AddTodo(ctx, "Original Title", "Original Description", model.PriorityLow, nil)

	err := svc.EditTodo(ctx, todo.ID, "Updated Title", "Updated Description", model.PriorityHigh, nil)
	if err != nil {
		t.Fatalf("failed to edit todo: %v", err)
	}

	// Verify changes
	updated, _ := repo.GetByID(ctx, todo.ID)
	if updated.Title != "Updated Title" {
		t.Errorf("expected title %q, got %q", "Updated Title", updated.Title)
	}
	if updated.Priority != model.PriorityHigh {
		t.Errorf("expected priority %d, got %d", model.PriorityHigh, updated.Priority)
	}
}

func TestTodoService_EditTodo_NotFound(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	err := svc.EditTodo(ctx, 999, "Title", "Description", model.PriorityMedium, nil)
	if err != ErrTodoNotFound {
		t.Errorf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestTodoService_DeleteTodo(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create a todo
	todo, _ := svc.AddTodo(ctx, "Test Todo", "", model.PriorityMedium, nil)

	// Delete it
	err := svc.DeleteTodo(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to delete todo: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, todo.ID)
	if err != repository.ErrTodoNotFound {
		t.Error("expected todo to be deleted")
	}
}

func TestTodoService_CompleteTodo(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create a pending todo
	todo, _ := svc.AddTodo(ctx, "Test Todo", "", model.PriorityMedium, nil)

	// Complete it
	err := svc.CompleteTodo(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to complete todo: %v", err)
	}

	// Verify status
	completed, _ := repo.GetByID(ctx, todo.ID)
	if completed.Status != model.StatusCompleted {
		t.Errorf("expected status %d, got %d", model.StatusCompleted, completed.Status)
	}
}

func TestTodoService_ListTodos(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create multiple todos
	svc.AddTodo(ctx, "Todo 1", "", model.PriorityLow, nil)
	svc.AddTodo(ctx, "Todo 2", "", model.PriorityMedium, nil)
	svc.AddTodo(ctx, "Todo 3", "", model.PriorityHigh, nil)

	todos, err := svc.ListTodos(ctx)
	if err != nil {
		t.Fatalf("failed to list todos: %v", err)
	}

	if len(todos) != 3 {
		t.Errorf("expected 3 todos, got %d", len(todos))
	}
}

func TestTodoService_ListPendingTodos(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create todos with different statuses
	todo1, _ := svc.AddTodo(ctx, "Pending 1", "", model.PriorityLow, nil)
	todo2, _ := svc.AddTodo(ctx, "Pending 2", "", model.PriorityMedium, nil)
	todo3, _ := svc.AddTodo(ctx, "To Complete", "", model.PriorityHigh, nil)

	// Complete one todo
	svc.CompleteTodo(ctx, todo3.ID)

	pending, err := svc.ListPendingTodos(ctx)
	if err != nil {
		t.Fatalf("failed to list pending todos: %v", err)
	}

	if len(pending) != 2 {
		t.Errorf("expected 2 pending todos, got %d", len(pending))
	}

	// Verify the correct todos are returned
	foundTodo1 := false
	foundTodo2 := false
	for _, todo := range pending {
		if todo.ID == todo1.ID {
			foundTodo1 = true
		}
		if todo.ID == todo2.ID {
			foundTodo2 = true
		}
	}

	if !foundTodo1 || !foundTodo2 {
		t.Error("expected to find both pending todos")
	}
}

func TestTodoService_ListCompletedTodos(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create todos
	svc.AddTodo(ctx, "Pending", "", model.PriorityLow, nil)
	todo2, _ := svc.AddTodo(ctx, "To Complete", "", model.PriorityMedium, nil)

	// Complete one
	svc.CompleteTodo(ctx, todo2.ID)

	completed, err := svc.ListCompletedTodos(ctx)
	if err != nil {
		t.Fatalf("failed to list completed todos: %v", err)
	}

	if len(completed) != 1 {
		t.Errorf("expected 1 completed todo, got %d", len(completed))
	}
}

func TestTodoService_ExportToJSON(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create some todos
	svc.AddTodo(ctx, "Todo 1", "Description 1", model.PriorityLow, nil)
	svc.AddTodo(ctx, "Todo 2", "Description 2", model.PriorityHigh, nil)

	// Export to temporary file
	tempDir := t.TempDir()
	exportPath := filepath.Join(tempDir, "export.json")

	err := svc.ExportToJSON(ctx, exportPath)
	if err != nil {
		t.Fatalf("failed to export: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		t.Error("export file was not created")
	}
}

func TestTodoService_ImportFromJSON(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create and export some todos
	svc.AddTodo(ctx, "Todo 1", "Description 1", model.PriorityLow, nil)
	svc.AddTodo(ctx, "Todo 2", "Description 2", model.PriorityHigh, nil)

	tempDir := t.TempDir()
	exportPath := filepath.Join(tempDir, "export.json")
	svc.ExportToJSON(ctx, exportPath)

	// Create a new service with empty repo
	newRepo := newMockRepository()
	newSvc := NewTodoService(newRepo)

	// Import the data
	err := newSvc.ImportFromJSON(ctx, exportPath)
	if err != nil {
		t.Fatalf("failed to import: %v", err)
	}

	// Verify imported data
	todos, _ := newSvc.ListTodos(ctx)
	if len(todos) != 2 {
		t.Errorf("expected 2 imported todos, got %d", len(todos))
	}
}

func TestTodoService_ImportFromJSON_FileNotFound(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	err := svc.ImportFromJSON(ctx, "/nonexistent/file.json")
	if err != ErrFileNotFound {
		t.Errorf("expected ErrFileNotFound, got %v", err)
	}
}

func TestTodoService_ImportFromJSON_InvalidJSON(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create a file with invalid JSON
	tempDir := t.TempDir()
	invalidPath := filepath.Join(tempDir, "invalid.json")
	os.WriteFile(invalidPath, []byte("invalid json content"), 0600)

	err := svc.ImportFromJSON(ctx, invalidPath)
	if err != ErrInvalidJSON {
		t.Errorf("expected ErrInvalidJSON, got %v", err)
	}
}

func TestTodoService_AddWorkDuration(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create a todo
	todo, err := svc.AddTodo(ctx, "Test Todo", "Description", model.PriorityMedium, nil)
	if err != nil {
		t.Fatalf("failed to add todo: %v", err)
	}

	// Add 25 minutes of work
	err = svc.AddWorkDuration(ctx, todo.ID, 25)
	if err != nil {
		t.Fatalf("failed to add work duration: %v", err)
	}

	// Verify work duration was added
	updated, err := repo.GetByID(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to get updated todo: %v", err)
	}

	if updated.WorkDuration != 25 {
		t.Errorf("expected work duration 25, got %d", updated.WorkDuration)
	}

	// Add another 25 minutes
	err = svc.AddWorkDuration(ctx, todo.ID, 25)
	if err != nil {
		t.Fatalf("failed to add work duration second time: %v", err)
	}

	// Verify cumulative work duration
	updated, err = repo.GetByID(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to get updated todo after second addition: %v", err)
	}

	if updated.WorkDuration != 50 {
		t.Errorf("expected cumulative work duration 50, got %d", updated.WorkDuration)
	}
}

func TestTodoService_AddWorkDuration_TodoNotFound(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Try to add work duration to non-existent todo
	err := svc.AddWorkDuration(ctx, 9999, 25)
	if err != ErrTodoNotFound {
		t.Errorf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestTodoService_AddWorkDuration_NegativeMinutes(t *testing.T) {
	repo := newMockRepository()
	svc := NewTodoService(repo)
	ctx := context.Background()

	// Create a todo
	todo, err := svc.AddTodo(ctx, "Test Todo", "Description", model.PriorityMedium, nil)
	if err != nil {
		t.Fatalf("failed to add todo: %v", err)
	}

	// Try to add negative minutes
	err = svc.AddWorkDuration(ctx, todo.ID, -10)
	if err != ErrInvalidWorkDuration {
		t.Errorf("expected ErrInvalidWorkDuration, got %v", err)
	}
}
