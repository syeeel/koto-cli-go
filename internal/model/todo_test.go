package model

import (
	"testing"
	"time"
)

func TestTodo_IsCompleted(t *testing.T) {
	tests := []struct {
		name   string
		status TodoStatus
		want   bool
	}{
		{
			name:   "completed todo returns true",
			status: StatusCompleted,
			want:   true,
		},
		{
			name:   "pending todo returns false",
			status: StatusPending,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := Todo{Status: tt.status}
			if got := todo.IsCompleted(); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodo_IsPending(t *testing.T) {
	tests := []struct {
		name   string
		status TodoStatus
		want   bool
	}{
		{
			name:   "pending todo returns true",
			status: StatusPending,
			want:   true,
		},
		{
			name:   "completed todo returns false",
			status: StatusCompleted,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := Todo{Status: tt.status}
			if got := todo.IsPending(); got != tt.want {
				t.Errorf("IsPending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodo_IsOverdue(t *testing.T) {
	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

	tests := []struct {
		name    string
		dueDate *time.Time
		status  TodoStatus
		want    bool
	}{
		{
			name:    "pending todo with past due date is overdue",
			dueDate: &past,
			status:  StatusPending,
			want:    true,
		},
		{
			name:    "pending todo with future due date is not overdue",
			dueDate: &future,
			status:  StatusPending,
			want:    false,
		},
		{
			name:    "completed todo with past due date is not overdue",
			dueDate: &past,
			status:  StatusCompleted,
			want:    false,
		},
		{
			name:    "todo without due date is not overdue",
			dueDate: nil,
			status:  StatusPending,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := Todo{
				DueDate: tt.dueDate,
				Status:  tt.status,
			}
			if got := todo.IsOverdue(); got != tt.want {
				t.Errorf("IsOverdue() = %v, want %v", got, tt.want)
			}
		})
	}
}
