package domain

import "github.com/google/uuid"

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
