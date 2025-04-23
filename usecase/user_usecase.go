package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/entity"
	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/repository"
)

// UserUseCase represents the user use case
type UserUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by their ID
func (uc *UserUseCase) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

// GetByEmail retrieves a user by their email
func (uc *UserUseCase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return uc.userRepo.GetByEmail(ctx, email)
}

// GetByUsername retrieves a user by their username
func (uc *UserUseCase) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return uc.userRepo.GetByUsername(ctx, username)
}

// Create creates a new user
func (uc *UserUseCase) Create(ctx context.Context, username, email, password, firstName, lastName string) (*entity.User, error) {
	// Create user entity
	user := entity.NewUser(username, email, password, firstName, lastName)

	// Validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	existingUser, err = uc.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Create user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates an existing user
func (uc *UserUseCase) Update(ctx context.Context, id uint64, username, email, firstName, lastName string) (*entity.User, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update user fields
	user.Username = username
	user.Email = email
	user.FirstName = firstName
	user.LastName = lastName
	user.UpdatedAt = time.Now()

	// Validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Update user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes a user by their ID
func (uc *UserUseCase) Delete(ctx context.Context, id uint64) error {
	return uc.userRepo.Delete(ctx, id)
}

// List retrieves a list of users with pagination
func (uc *UserUseCase) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return uc.userRepo.List(ctx, limit, offset)
}
