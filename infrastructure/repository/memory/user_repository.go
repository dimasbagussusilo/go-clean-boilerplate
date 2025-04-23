package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/entity"
	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/repository"
)

// Ensure UserRepository implements repository.UserRepository
var _ repository.UserRepository = (*UserRepository)(nil)

// UserRepository is an in-memory implementation of repository.UserRepository
type UserRepository struct {
	mu    sync.RWMutex
	users map[uint64]*entity.User
	// Auto-increment ID
	lastID uint64
}

// NewUserRepository creates a new in-memory user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[uint64]*entity.User),
		lastID: 0,
	}
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetByUsername retrieves a user by their username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return errors.New("email already exists")
		}
		if existingUser.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	// Assign ID
	r.lastID++
	user.ID = r.lastID

	// Store user
	r.users[user.ID] = user

	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Check if email already exists for another user
	for id, existingUser := range r.users {
		if id != user.ID && existingUser.Email == user.Email {
			return errors.New("email already exists")
		}
		if id != user.ID && existingUser.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	// Update user
	r.users[user.ID] = user

	return nil
}

// Delete deletes a user by their ID
func (r *UserRepository) Delete(ctx context.Context, id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)

	return nil
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	// Apply pagination
	if offset >= len(users) {
		return []*entity.User{}, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], nil
}
