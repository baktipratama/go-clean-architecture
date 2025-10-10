package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user with validation
func NewUser(name, email string) (*User, error) {
	if err := validateUserInput(name, email); err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateName updates the user's name with validation
func (u *User) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail updates the user's email with validation
func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	// Basic email validation - in production you'd want more robust validation
	if !isValidEmail(email) {
		return ErrInvalidEmail
	}
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// validateUserInput validates the input for creating a user
func validateUserInput(name, email string) error {
	if name == "" {
		return ErrInvalidName
	}
	if email == "" {
		return ErrInvalidEmail
	}
	if !isValidEmail(email) {
		return ErrInvalidEmail
	}
	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	// Basic email validation - in production use a proper email validation library
	if len(email) < 5 {
		return false
	}

	atCount := 0
	dotAfterAt := false
	atIndex := -1

	for i, char := range email {
		if char == '@' {
			atCount++
			atIndex = i
		}
		if char == '.' && i > atIndex && atIndex != -1 {
			dotAfterAt = true
		}
	}

	return atCount == 1 && dotAfterAt && atIndex > 0 && atIndex < len(email)-1
}
