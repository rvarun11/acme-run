package amqphandler

import (
	"time"

	"github.com/google/uuid"
)

type ShelterAvailable struct {
	// WorkoutID for which the Shelter Availability is there or not
	WorkoutID uuid.UUID `json:"workout_id"`
	// ShelterAvailable or not
	ShelterAvailable bool `json:"shelter_available"`
	// Distance to Shelter
	DistanceToShelter float64 `json:"distance_to_shelter"`
}

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
