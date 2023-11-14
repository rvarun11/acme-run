package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserName  = errors.New("a user has to have a name")
	ErrInvalidUserEmail = errors.New("a user has to have a email")
	ErrInvalidUserDOB   = errors.New("a user has to have an DOB")
)

type User struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID
	// Name of the user
	Name string
	// Email
	Email string
	// Date of Birth of the player TODO: Fix type
	DateOfBirth string
}

// NewPlayer is a factory to create a new Player aggregate
func NewUser(name string, email string, dob string) (*User, error) {
	// Validate that the Email has @, TODO: add more validation
	if name == "" {
		return &User{}, ErrInvalidUserName
	}
	if email == "" {
		return &User{}, ErrInvalidUserEmail
	}
	if dob == "" {
		return &User{}, ErrInvalidUserDOB
	}

	// Create a new user and generate ID
	user := &User{
		ID:          uuid.New(),
		Name:        name,
		Email:       email,
		DateOfBirth: dob,
	}

	return user, nil
}
