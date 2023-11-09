package httphandler

import "github.com/google/uuid"

type ShelterAvailable struct {
	// WorkoutID for which the Shelter Availability is there or not
	WorkoutID uuid.UUID `json:"workout_id"`
	// ShelterAvailable or not
	ShelterAvailable bool `json:"shelter_available"`
	// Distance to Shelter
	DistanceToShelter float64 `json:"distance_to_shelter"`
}
