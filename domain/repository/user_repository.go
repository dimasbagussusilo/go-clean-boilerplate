package repository

import (
	"context"

	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/entity"
)

// UserRepository represents the user repository contract
type UserRepository interface {
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uint64) (*entity.User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes a user by their ID
	Delete(ctx context.Context, id uint64) error

	// List retrieves a list of users with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}
