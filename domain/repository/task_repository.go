package repository

import (
	"context"

	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/entity"
)

// TaskRepository represents the task repository contract
type TaskRepository interface {
	// GetByID retrieves a task by its ID
	GetByID(ctx context.Context, id uint64) (*entity.Task, error)

	// GetByUserID retrieves tasks by user ID
	GetByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*entity.Task, error)

	// Create creates a new task
	Create(ctx context.Context, task *entity.Task) error

	// Update updates an existing task
	Update(ctx context.Context, task *entity.Task) error

	// Delete deletes a task by its ID
	Delete(ctx context.Context, id uint64) error

	// List retrieves a list of tasks with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.Task, error)
}
