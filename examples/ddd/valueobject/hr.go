// Package entities holds all the entities that are shared across all subdomains
package entity

import (
	"time"

	"github.com/google/uuid"
)

// HeartRate is a entity that represents a heart in all Domains
type HeartRate struct {
	// workout id is the id of the workout
	workoutId uuid.UUID
	// heartRate is the recorded heart rate by the HRM
	heartRate float64
	// createdAt is the time when the heartRate was recorded
	createdAt time.Time
}
