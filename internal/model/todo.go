package model

import (
	"fmt"
	"time"
)

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
	ID           int64      `db:"id"`
	Title        string     `db:"title"`
	Description  string     `db:"description"`
	Status       TodoStatus `db:"status"`
	Priority     Priority   `db:"priority"`
	DueDate      *time.Time `db:"due_date"`
	WorkDuration int        `db:"work_duration"` // Cumulative work time in minutes
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
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

// GetWorkDurationFormatted returns the work duration in human-readable format
// Examples:
//   - 0 minutes: ""
//   - 25 minutes: "25m"
//   - 60 minutes: "1h"
//   - 125 minutes: "2h 5m"
func (t Todo) GetWorkDurationFormatted() string {
	// Defensive: handle negative or zero duration
	if t.WorkDuration <= 0 {
		return ""
	}

	hours := t.WorkDuration / 60
	minutes := t.WorkDuration % 60

	// Hours only (no remaining minutes)
	if hours > 0 && minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	// Minutes only (less than an hour)
	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	// Hours and minutes
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
