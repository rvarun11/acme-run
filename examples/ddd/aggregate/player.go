// Package entities holds all the entities that are shared across all subdomains
package aggregate

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/entity"
)

var (
	ErrInvalidEmail = errors.New("a customer has to have a valid email")
)

// Player is a entity that represents a Player in all Domains
type Player struct {
	// User is the root entity of Player
	user *entity.User

	weight float32
	// Height of the player
	height float32
	// GeographicalZone is a group of trails in a region
	geographicalzone string
	// CreatedAt is the time when the player registered
	createdat time.Time
	// UpdatedAt is the time when the player last updated the profile
	updatedat time.Time
}

// NewPlayer is a factory to create a new Player aggregate
// It will validate that the name is not empty
func NewPlayer(name string, email string, dob string, weight float32, height float32, geographicalzone string) (Player, error) {
	// Validate that the Email has @
	_, err := mail.ParseAddress(email)
	if err != nil {
		return Player{}, ErrInvalidEmail
	}

	// Create a new person and generate ID
	user := &entity.User{
		Name:        name,
		Email:       email,
		DateOfBirth: dob,
		ID:          uuid.New(),
	}
	// Create a customer object and initialize all the values to avoid nil pointer exceptions
	return Player{
		user:             user,
		weight:           weight,
		height:           height,
		geographicalzone: geographicalzone,
	}, nil
}

func (player *Player) GetID() uuid.UUID {
	return player.user.ID
}

func (player *Player) SetID(id uuid.UUID) {
	player.user.ID = id
}

func (player *Player) GetEmail() string {
	return player.user.Email
}

func (player *Player) SetEmail(email string) {
	player.user.Email = email
}
