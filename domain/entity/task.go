package entity

import (
	"time"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	// TaskStatusPending represents a pending task
	TaskStatusPending TaskStatus = "pending"
	// TaskStatusInProgress represents a task in progress
	TaskStatusInProgress TaskStatus = "in_progress"
	// TaskStatusCompleted represents a completed task
	TaskStatusCompleted TaskStatus = "completed"
)

// Task represents the task entity
type Task struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	UserID      uint64     `json:"user_id"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NewTask creates a new task
func NewTask(title, description string, userID uint64, dueDate *time.Time) *Task {
	now := time.Now()
	return &Task{
		Title:       title,
		Description: description,
		Status:      TaskStatusPending,
		UserID:      userID,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate validates the task entity
func (t *Task) Validate() error {
	// TODO: Implement validation logic
	return nil
}

// MarkInProgress marks the task as in progress
func (t *Task) MarkInProgress() {
	t.Status = TaskStatusInProgress
	t.UpdatedAt = time.Now()
}

// MarkCompleted marks the task as completed
func (t *Task) MarkCompleted() {
	t.Status = TaskStatusCompleted
	t.UpdatedAt = time.Now()
}
