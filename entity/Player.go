// Package entities holds all the entities that are shared across all subdomains
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Person is a entity that represents a person in all Domains
type Player struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID
	// First Name of the person
	FirstName string
	// Last Name of the person
	LastName string
	// Email
	Email string
	// Date of Birth of the player TODO: Fix type
	DateOfBirth string
	// Weight of the player
	Weight float32
	// Height of the player
	Height float32
	// GeographicalZone is a group of trails in a region
	GeographicalZone string
	// CreatedAt is the time when the player registered
	CreatedAt time.Time
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time
}
