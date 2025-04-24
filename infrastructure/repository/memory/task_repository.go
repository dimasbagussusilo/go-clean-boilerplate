package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/entity"
	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/repository"
)

// Ensure TaskRepository implements repository.TaskRepository
var _ repository.TaskRepository = (*TaskRepository)(nil)

// TaskRepository is an in-memory implementation of repository.TaskRepository
type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[uint64]*entity.Task
	// Auto-increment ID
	lastID uint64
}

// NewTaskRepository creates a new in-memory task repository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks:  make(map[uint64]*entity.Task),
		lastID: 0,
	}
}

// GetByID retrieves a task by its ID
func (r *TaskRepository) GetByID(_ context.Context, id uint64) (*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

// GetByUserID retrieves tasks by user ID
func (r *TaskRepository) GetByUserID(_ context.Context, userID uint64, limit, offset int) ([]*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Filter tasks by user ID
	userTasks := make([]*entity.Task, 0)
	for _, task := range r.tasks {
		if task.UserID == userID {
			userTasks = append(userTasks, task)
		}
	}

	// Apply pagination
	if offset >= len(userTasks) {
		return []*entity.Task{}, nil
	}

	end := offset + limit
	if end > len(userTasks) {
		end = len(userTasks)
	}

	return userTasks[offset:end], nil
}

// Create creates a new task
func (r *TaskRepository) Create(_ context.Context, task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Assign ID
	r.lastID++
	task.ID = r.lastID

	// Store task
	r.tasks[task.ID] = task

	return nil
}

// Update updates an existing task
func (r *TaskRepository) Update(_ context.Context, task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}

	// Update task
	r.tasks[task.ID] = task

	return nil
}

// Delete deletes a task by its ID
func (r *TaskRepository) Delete(_ context.Context, id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, id)

	return nil
}

// List retrieves a list of tasks with pagination
func (r *TaskRepository) List(_ context.Context, limit, offset int) ([]*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	tasks := make([]*entity.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	// Apply pagination
	if offset >= len(tasks) {
		return []*entity.Task{}, nil
	}

	end := offset + limit
	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[offset:end], nil
}
