// Package entities holds all the entities that are shared across all subdomains
package entity

import (
	"github.com/google/uuid"
)

// Person is a entity that represents a person in all Domains
type HeartRateMonitor struct {
	// WorkoutSessionID is the .....
	WorkoutSessionID uuid.UUID
	//
	Min_HR float32
	//
	Max_HR float32
	// HardcoreMode is the difficulty level chosen by the player
	Interval float32
}
