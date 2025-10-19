package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver

	"github.com/syeeel/koto-cli-go/internal/model"
)

// schemaSQL contains the database schema initialization SQL
const schemaSQL = `
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    priority INTEGER NOT NULL DEFAULT 0,
    due_date DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
`

var (
	// ErrTodoNotFound is returned when a todo is not found
	ErrTodoNotFound = errors.New("todo not found")
)

// SQLiteRepository implements TodoRepository using SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// Open database connection using modernc.org/sqlite (Pure Go, no CGO required)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set file permissions for security (only if not in-memory database)
	if dbPath != ":memory:" {
		if err := os.Chmod(dbPath, 0600); err != nil {
			// Ignore error if file doesn't exist yet
			// (it will be created by SQLite)
		}
	}

	// Initialize schema
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &SQLiteRepository{db: db}, nil
}

// initSchema initializes the database schema
func initSchema(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	return err
}

// Create creates a new todo item
func (r *SQLiteRepository) Create(ctx context.Context, todo *model.Todo) error {
	query := `
		INSERT INTO todos (title, description, status, priority, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.Priority,
		todo.DueDate,
		todo.CreatedAt,
		todo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	todo.ID = id
	return nil
}

// GetByID retrieves a todo by ID
func (r *SQLiteRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM todos
		WHERE id = ?
	`

	todo := &model.Todo{}
	var dueDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.Priority,
		&dueDate,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTodoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	if dueDate.Valid {
		todo.DueDate = &dueDate.Time
	}

	return todo, nil
}

// GetAll retrieves all todos
func (r *SQLiteRepository) GetAll(ctx context.Context) ([]*model.Todo, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

// GetByStatus retrieves todos by status
func (r *SQLiteRepository) GetByStatus(ctx context.Context, status model.TodoStatus) ([]*model.Todo, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM todos
		WHERE status = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos by status: %w", err)
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

// Update updates a todo
func (r *SQLiteRepository) Update(ctx context.Context, todo *model.Todo) error {
	query := `
		UPDATE todos
		SET title = ?, description = ?, status = ?, priority = ?, due_date = ?, updated_at = ?
		WHERE id = ?
	`

	todo.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.Priority,
		todo.DueDate,
		todo.UpdatedAt,
		todo.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}

// Delete deletes a todo by ID
func (r *SQLiteRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM todos WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}

// MarkAsCompleted marks a todo as completed
func (r *SQLiteRepository) MarkAsCompleted(ctx context.Context, id int64) error {
	query := `
		UPDATE todos
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, model.StatusCompleted, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark todo as completed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}

// Close closes the repository connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// scanTodos is a helper function to scan multiple todo rows
func (r *SQLiteRepository) scanTodos(rows *sql.Rows) ([]*model.Todo, error) {
	var todos []*model.Todo

	for rows.Next() {
		todo := &model.Todo{}
		var dueDate sql.NullTime

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&todo.Priority,
			&dueDate,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}

		if dueDate.Valid {
			todo.DueDate = &dueDate.Time
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todos: %w", err)
	}

	return todos, nil
}
