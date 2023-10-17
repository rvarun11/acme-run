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
	WorkoutId uuid.UUID `json:"workout_id"`
	HRMId     uuid.UUID `json:"hrm_id"`
	HRate     string    `json:"heart_rate"`
	CreatedAt time.Time `json:"created_at"`
}

// Getters and Setters for HRM
func (hrm *HRM) GetID() uuid.UUID {
	return hrm.HRMId
}

func (hrm *HRM) SetID(id uuid.UUID) {
	hrm.HRMId = id
}

/*
	func (hrm *HRM) getState() string {
		return hrm.Status
	}

	func (hrm *HRM) connectToHRM() {
		hrm.Status = "Connected"
	}
*/
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
	}
	return hrmN, nil
}
