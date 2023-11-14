package domain

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail = errors.New("a customer has to have a valid email")
)

type Preference string

const (
	Strength Preference = "strength"
	Cardio   Preference = "cardio"
)

// Player is a entity that represents a Player in all Domains
type Player struct {
	// User is the root entity of player
	ID uuid.UUID
	// User is the root entity of player
	User *User
	// Weight of the player
	Weight float64
	// Height of the player
	Height float64
	// Preference of player
	Preference Preference
	// GeographicalZone is a group of trails in a region
	ZoneID uuid.UUID
	// CreatedAt is the time when the player registered
	CreatedAt time.Time
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time
}

// NewPlayer is a factory to create a new Player aggregate
func NewPlayer(name string, email string, dob string, weight float64, height float64, pref Preference, zoneID uuid.UUID) (*Player, error) {
	// Validate that the Email has @, TODO: add more validation
	_, err := mail.ParseAddress(email)
	if err != nil {
		return &Player{}, ErrInvalidEmail
	}

	// Create a new user and generate ID
	user := User{
		ID:          uuid.New(),
		Name:        name,
		Email:       email,
		DateOfBirth: dob,
	}
	// Create a customer object and initialize all the values to avoid nil pointer exceptions
	player := &Player{
		ID:         uuid.New(),
		User:       &user,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Weight:     weight,
		Height:     height,
		Preference: pref,
		ZoneID:     zoneID, // TODO: This is a temp field
	}

	return player, nil
}
