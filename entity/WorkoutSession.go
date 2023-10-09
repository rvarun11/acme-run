// Package entities holds all the entities that are shared across all subdomains
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Person is a entity that represents a person in all Domains
type Workout struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID
	// Name of the Workout
	Name string
	// CreatedAt is the time when the workout was started/created at?
	CreatedAt time.Time
	// TODO: Duration of the workout
	Duration float64

	// DurationCovered is the total distance covered during the session
	DistanceCovered float64
	// workouttype can be either Physical Workout or Cardio
	WorkoutType string
	// HardcoreMode is the difficulty level chosen by the player
	HardcoreMode bool
	// Avergage HRM Reading from the workout TODO: See if it belongs here
	AverageHeartRates float64
}
