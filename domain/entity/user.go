package entity

import (
	"time"
)

// User represents the user entity
type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is not exposed in JSON
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user
func NewUser(username, email, password, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Email:     email,
		Password:  password, // Note: In a real application, this should be hashed
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the user entity
func (u *User) Validate() error {
	// TODO: Implement validation logic
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}