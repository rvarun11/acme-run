package domain

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

var (
	ErrInvalidEmail = errors.New("a customer has to have a valid email")
)

type HRM struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	HRMId uuid.UUID `json:"hrmid"`
	// Name of the user
	Status string `json:"status"`
}

// Player is a entity that represents a Player in all Domains
type HRM struct {
	// Email
	HRate     string    `json:"hrate"`
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time `json:"updated_at"`
}

// Getters and Setters for Player
func (hrm *HRM) GetID() uuid.UUID {
	return hrm.HRMId
}

func (hrm *HRM) SetID(id uuid.UUID) {
	hrm.HRMId = id
}

func (hrm *HRM) getHRate() string {
	return hrm.HRate
}

func (hrm *HRM) readHRate() {
	hrm.HRate = "100"
}

// NewPlayer is a factory to create a new Player aggregate

func NewHRM(hrm HRM) (HRM, error) {

	// Create a hrm object and initialize all the values to avoid nil pointer exceptions
	hrmN := HRM{
		HRMId:     hrm.HRMId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "connected"
	}
	return hrmN, nil
}
