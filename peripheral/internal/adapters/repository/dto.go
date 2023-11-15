package repository

import (
	"time"

	"github.com/google/uuid"
)

type LastLocation struct {
	// WorkoutID for which the Shelter Availability is there or not
	WorkoutID uuid.UUID `json:"workout_id"`
	// Latitude of the Player
	Latitude float64 `json:"latitude"`
	// Longitude of the Player
	Longitude float64 `json:"longitude"`
	// Time of location
	TimeOfLocation time.Time `json:"time_of_location"`
}

type LastHR struct {
	HRMID uuid.UUID `json:"hrm_id"`
	// Latitude of the Player
	HeartRate int `json:"heart_rate"`
	// Time of reading
	TimeOfLocation time.Time `json:"time_of_location"`
}
