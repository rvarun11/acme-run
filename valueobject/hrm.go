// Package entities holds all the entities that are shared across all subdomains
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Person is a entity that represents a person in all Domains
type HeartRateMonitor struct {
	// WorkoutID is the .....
	WorkoutID uuid.UUID
	//
	minHeartRate float32
	//
	maxHeartRate float32
	// HardcoreMode is the difficulty level chosen by the player
	interval float32

	createdAt time.Time
}
