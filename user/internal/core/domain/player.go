package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPlayerHeight     = errors.New("a player has to have a valid height")
	ErrInvalidPlayerWeight     = errors.New("a player has to have a valid weight")
	ErrInvalidPlayerPreference = errors.New("a player has to have a valid preference")
	ErrInvalidZoneID           = errors.New("a player must belong to a valid zone")
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

	user, err := NewUser(name, email, dob)
	if err != nil {
		return &Player{}, err
	}

	err = validateHeight(height)
	if err != nil {
		return &Player{}, err
	}

	err = validateWeight(weight)
	if err != nil {
		return &Player{}, err
	}

	err = validatePreference(pref)
	if err != nil {
		return &Player{}, err
	}

	err = validateZoneID(zoneID)
	if err != nil {
		return &Player{}, err
	}

	player := &Player{
		ID:         uuid.New(),
		User:       user,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Weight:     weight,
		Height:     height,
		Preference: pref,
		ZoneID:     zoneID,
	}

	return player, nil
}

func validateHeight(h float64) error {
	if h == 0.0 {
		return ErrInvalidPlayerHeight
	}
	return nil
}

func validateWeight(w float64) error {
	if w == 0.0 {
		return ErrInvalidPlayerWeight
	}
	return nil
}

func validatePreference(p Preference) error {
	switch p {
	case Strength, Cardio:
		return nil
	default:
		return ErrInvalidPlayerPreference
	}
}

func validateZoneID(zid uuid.UUID) error {
	if zid == uuid.Nil {
		return ErrInvalidZoneID
	}
	return nil
}
