// Package entities holds all the entities that are shared across all subdomains
package entity

import (
	"github.com/google/uuid"
)

// User is a entity that represents a Users in all Domains
type User struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID
	// First Name of the person
	Name string
	// Email
	Email string
	// Date of Birth of the player TODO: Fix type
	DateOfBirth string
}
