package repository

import (
	"context"
	"testing"
	"time"

	"github.com/syeeel/koto-cli-go/internal/model"
)

// setupTestDB creates a new in-memory SQLite database for testing
func setupTestDB(t *testing.T) *SQLiteRepository {
	t.Helper()

	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create test repository: %v", err)
	}

	return repo
}

func TestNewSQLiteRepository(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	if repo == nil {
		t.Fatal("expected repository to be created, got nil")
	}

	if repo.db == nil {
		t.Fatal("expected database connection to be initialized, got nil")
	}
}

func TestSQLiteRepository_Create(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()
	now := time.Now()
	dueDate := now.Add(24 * time.Hour)

	todo := &model.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Status:      model.StatusPending,
		Priority:    model.PriorityMedium,
		DueDate:     &dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := repo.Create(ctx, todo)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	if todo.ID == 0 {
		t.Error("expected todo ID to be set after creation")
	}
}

func TestSQLiteRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()
	now := time.Now()

	// Create a todo
	todo := &model.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Status:      model.StatusPending,
		Priority:    model.PriorityHigh,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := repo.Create(ctx, todo); err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// Retrieve the todo
	retrieved, err := repo.GetByID(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to get todo: %v", err)
	}

	if retrieved.ID != todo.ID {
		t.Errorf("expected ID %d, got %d", todo.ID, retrieved.ID)
	}
	if retrieved.Title != todo.Title {
		t.Errorf("expected title %q, got %q", todo.Title, retrieved.Title)
	}
	if retrieved.Description != todo.Description {
		t.Errorf("expected description %q, got %q", todo.Description, retrieved.Description)
	}
	if retrieved.Status != todo.Status {
		t.Errorf("expected status %d, got %d", todo.Status, retrieved.Status)
	}
	if retrieved.Priority != todo.Priority {
		t.Errorf("expected priority %d, got %d", todo.Priority, retrieved.Priority)
	}
}

func TestSQLiteRepository_GetByID_NotFound(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	if err != ErrTodoNotFound {
		t.Errorf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestSQLiteRepository_GetAll(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()
	now := time.Now()

	// Create multiple todos
	todos := []*model.Todo{
		{Title: "Todo 1", Status: model.StatusPending, Priority: model.PriorityLow, CreatedAt: now, UpdatedAt: now},
		{Title: "Todo 2", Status: model.StatusCompleted, Priority: model.PriorityMedium, CreatedAt: now, UpdatedAt: now},
		{Title: "Todo 3", Status: model.StatusPending, Priority: model.PriorityHigh, CreatedAt: now, UpdatedAt: now},
	}

	for _, todo := range todos {
		if err := repo.Create(ctx, todo); err != nil {
			t.Fatalf("failed to create todo: %v", err)
		}
	}

	// Retrieve all todos
	allTodos, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all todos: %v", err)
	}

	if len(allTodos) != 3 {
		t.Errorf("expected 3 todos, got %d", len(allTodos))
	}
}

func TestSQLiteRepository_GetByStatus(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()
	now := time.Now()

	// Create todos with different statuses
	todos := []*model.Todo{
		{Title: "Pending 1", Status: model.StatusPending, Priority: model.PriorityLow, CreatedAt: now, UpdatedAt: now},
		{Title: "Completed 1", Status: model.StatusCompleted, Priority: model.PriorityMedium, CreatedAt: now, UpdatedAt: now},
		{Title: "Pending 2", Status: model.StatusPending, Priority: model.PriorityHigh, CreatedAt: now, UpdatedAt: now},
	}

	for _, todo := range todos {
		if err := repo.Create(ctx, todo); err != nil {
			t.Fatalf("failed to create todo: %v", err)
		}
	}

	// Get pending todos
	pendingTodos, err := repo.GetByStatus(ctx, model.StatusPending)
	if err != nil {
		t.Fatalf("failed to get pending todos: %v", err)
	}

	if len(pendingTodos) != 2 {
		t.Errorf("expected 2 pending todos, got %d", len(pendingTodos))
	}

	// Get completed todos
	completedTodos, err := repo.GetByStatus(ctx, model.StatusCompleted)
	if err != nil {
		t.Fatalf("failed to get completed todos: %v", err)
	}

	if len(completedTodos) != 1 {
		t.Errorf("expected 1 completed todo, got %d", len(completedTodos))
	}
}

func TestSQLiteRepository_Update(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()
	now := time.Now()

	// Create a todo
	todo := &model.Todo{
		Title:       "Original Title",
		Description: "Original Description",
		Status:      model.StatusPending,
		Priority:    model.PriorityLow,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := repo.Create(ctx, todo); err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// Update the todo
	todo.Title = "Updated Title"
	todo.Description = "Updated Description"
	todo.Priority = model.PriorityHigh

	if err := repo.Update(ctx, todo); err != nil {
		t.Fatalf("failed to update todo: %v", err)
	}

	// Retrieve and verify
	updated, err := repo.GetByID(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to get updated todo: %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("expected title %q, got %q", "Updated Title", updated.Title)
	}
	if updated.Description != "Updated Description" {
		t.Errorf("expected description %q, got %q", "Updated Description", updated.Description)
	}
	if updated.Priority != model.PriorityHigh {
		t.Errorf("expected priority %d, got %d", model.PriorityHigh, updated.Priority)
	}
}

func TestSQLiteRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()
	now := time.Now()

	// Create a todo
	todo := &model.Todo{
		Title:     "Todo to Delete",
		Status:    model.StatusPending,
		Priority:  model.PriorityMedium,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Create(ctx, todo); err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// Delete the todo
	if err := repo.Delete(ctx, todo.ID); err != nil {
		t.Fatalf("failed to delete todo: %v", err)
	}

	// Verify deletion
	_, err := repo.GetByID(ctx, todo.ID)
	if err != ErrTodoNotFound {
		t.Errorf("expected ErrTodoNotFound after deletion, got %v", err)
	}
}

func TestSQLiteRepository_MarkAsCompleted(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	ctx := context.Background()
	now := time.Now()

	// Create a pending todo
	todo := &model.Todo{
		Title:     "Pending Todo",
		Status:    model.StatusPending,
		Priority:  model.PriorityMedium,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Create(ctx, todo); err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// Mark as completed
	if err := repo.MarkAsCompleted(ctx, todo.ID); err != nil {
		t.Fatalf("failed to mark todo as completed: %v", err)
	}

	// Verify status
	completed, err := repo.GetByID(ctx, todo.ID)
	if err != nil {
		t.Fatalf("failed to get completed todo: %v", err)
	}

	if completed.Status != model.StatusCompleted {
		t.Errorf("expected status %d, got %d", model.StatusCompleted, completed.Status)
	}
}
