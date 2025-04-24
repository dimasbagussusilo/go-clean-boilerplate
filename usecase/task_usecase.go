package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/entity"
	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/repository"
)

// TaskUseCase represents the task use case
type TaskUseCase struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
}

// NewTaskUseCase creates a new task use case
func NewTaskUseCase(taskRepo repository.TaskRepository, userRepo repository.UserRepository) *TaskUseCase {
	return &TaskUseCase{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

// GetByID retrieves a task by its ID
func (uc *TaskUseCase) GetByID(ctx context.Context, id uint64) (*entity.Task, error) {
	return uc.taskRepo.GetByID(ctx, id)
}

// GetByUserID retrieves tasks by user ID
func (uc *TaskUseCase) GetByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*entity.Task, error) {
	// Verify user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return uc.taskRepo.GetByUserID(ctx, userID, limit, offset)
}

// Create creates a new task
func (uc *TaskUseCase) Create(ctx context.Context, title, description string, userID uint64, dueDate *time.Time) (*entity.Task, error) {
	// Verify user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Create task entity
	task := entity.NewTask(title, description, userID, dueDate)

	// Validate task
	if err := task.Validate(); err != nil {
		return nil, err
	}

	// Create task
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Update updates an existing task
func (uc *TaskUseCase) Update(ctx context.Context, id uint64, title, description string, status entity.TaskStatus, dueDate *time.Time) (*entity.Task, error) {
	// Get existing task
	task, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update task fields
	task.Title = title
	task.Description = description
	task.Status = status
	task.DueDate = dueDate
	task.UpdatedAt = time.Now()

	// Validate task
	if err := task.Validate(); err != nil {
		return nil, err
	}

	// Update task
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Delete deletes a task by its ID
func (uc *TaskUseCase) Delete(ctx context.Context, id uint64) error {
	return uc.taskRepo.Delete(ctx, id)
}

// List retrieves a list of tasks with pagination
func (uc *TaskUseCase) List(ctx context.Context, limit, offset int) ([]*entity.Task, error) {
	return uc.taskRepo.List(ctx, limit, offset)
}

// MarkInProgress marks a task as in progress
func (uc *TaskUseCase) MarkInProgress(ctx context.Context, id uint64) (*entity.Task, error) {
	// Get existing task
	task, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Mark as in progress
	task.MarkInProgress()

	// Update task
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// MarkCompleted marks a task as completed
func (uc *TaskUseCase) MarkCompleted(ctx context.Context, id uint64) (*entity.Task, error) {
	// Get existing task
	task, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Mark as completed
	task.MarkCompleted()

	// Update task
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}