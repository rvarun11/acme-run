package dto

import "github.com/google/uuid"


type UserDTO struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID `json:"id"`
	// Name of the user
	Name string `json:"name"`
	// Email
	Email string `json:"email"`
	// DoB
	DateOfBirth string `json:"dob"`
}
