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

type User struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID `json:"id"`
	// Name of the user
	Name string `json:"name"`
	// Email
	Email string `json:"email"`
	// Date of Birth of the player TODO: Fix type
	DateOfBirth string `json:"dob"`
}

// Player is a entity that represents a Player in all Domains
type Player struct {
	// User is the root entity of player
	User *User `json:"user"`
	// Weight of the player
	Weight float64 `json:"weight"`
	// Height of the player
	Height float64 `json:"height"`
	// GeographicalZone is a group of trails in a region
	GeographicalZone string `json:"geographical_zone"`
	// CreatedAt is the time when the player registered
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time `json:"updated_at"`
}

// Getters and Setters for Player
func (player *Player) GetID() uuid.UUID {
	return player.User.ID
}

func (player *Player) SetID(id uuid.UUID) {
	player.User.ID = id
}

// NewPlayer is a factory to create a new Player aggregate
func NewPlayer(p Player) (Player, error) {
	// Validate that the Email has @
	_, err := mail.ParseAddress(p.User.Email)
	if err != nil {
		return Player{}, ErrInvalidEmail
	}

	// Create a new user and generate ID
	user := &User{
		Name:        p.User.Name,
		Email:       p.User.Email,
		DateOfBirth: p.User.DateOfBirth,
		ID:          uuid.New(),
	}
	// Create a customer object and initialize all the values to avoid nil pointer exceptions
	player := Player{
		User:             user,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Weight:           p.Weight,
		Height:           p.Height,
		GeographicalZone: p.GeographicalZone,
	}

	return player, nil
}
