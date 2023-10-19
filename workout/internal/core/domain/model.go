// Package entities holds all the entities that are shared across all subdomains
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidWorkout = errors.New("no playerId or trailId associated with the workout session")
)

// Workout is a entity that represents a workout in all Domains
type Workout struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID `json:"id"`
	// trailId is the id of the trail player is on
	TrailID uuid.UUID `json:"trail_id"`
	// PlayerID of the player starting the workout session
	PlayerID uuid.UUID `json:"player_id"`
	// InProgress tells whether the workout is in progress
	IsCompleted bool `json:"is_completed"`
	// CreatedAt is the time when the workout was started/created at?
	CreatedAt time.Time `json:"created_at"`
	// Duration of the workout, TODO: fix type
	EndedAt time.Time `json:"ended_at"`
	// DurationCovered is the total distance covered during the session
	DistanceCovered float64 `json:"distance_covered"`
	// TODO: temp value. It can be either "cardio", "physical" or "dynamic"
	Category string `json:"category"`
	// HardcoreMode is the difficulty level chosen by the player
	HardcoreMode bool `json:"HardcoreMode"`
	//  HRM Reading from the workout
	// TODO: HeartRate should be a valueobject of hrmValue + created_at
	HeartRate []uint16
}

func NewWorkout(w Workout) (Workout, error) {
	if w.PlayerID == uuid.Nil || w.TrailID == uuid.Nil {
		return Workout{}, ErrInvalidWorkout
	}

	return Workout{
		ID:              uuid.New(),
		PlayerID:        w.PlayerID,
		TrailID:         w.TrailID,
		Category:        w.Category,
		IsCompleted:     false,
		HardcoreMode:    w.HardcoreMode,
		CreatedAt:       time.Now(),
		EndedAt:         time.Time{},
		DistanceCovered: 0,
		HeartRate:       []uint16{},
		// heartRates:      make([]valueobject.HeartRate, 0),
	}, nil
}

func (w *Workout) GetID() uuid.UUID {
	return w.ID
}

func (w *Workout) SetID(id uuid.UUID) {
	w.ID = id
}
