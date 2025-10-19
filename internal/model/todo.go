package model

import "time"

// TodoStatus represents the status of a todo item
type TodoStatus int

const (
	// StatusPending indicates a pending todo
	StatusPending TodoStatus = iota
	// StatusCompleted indicates a completed todo
	StatusCompleted
)

// Priority represents the priority level of a todo item
type Priority int

const (
	// PriorityLow indicates low priority
	PriorityLow Priority = iota
	// PriorityMedium indicates medium priority
	PriorityMedium
	// PriorityHigh indicates high priority
	PriorityHigh
)

// Todo represents a todo item
type Todo struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      TodoStatus `db:"status"`
	Priority    Priority   `db:"priority"`
	DueDate     *time.Time `db:"due_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// IsCompleted returns true if the todo is completed
func (t Todo) IsCompleted() bool {
	return t.Status == StatusCompleted
}

// IsPending returns true if the todo is pending
func (t Todo) IsPending() bool {
	return t.Status == StatusPending
}

// IsOverdue returns true if the todo is overdue (past due date and still pending)
func (t Todo) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && t.IsPending()
}
