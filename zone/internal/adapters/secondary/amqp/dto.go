package amqp

import "github.com/google/uuid"

type ShelterDTO struct {
	WorkoutID           uuid.UUID `json:"workout_id"`
	ShelterAvailability bool      `json:"shelter_availability"`
	// Shelter ID
	ShelterID uuid.UUID `json:"shelter_id"`
	// Shelter name
	ShelterName string `json:"shelter_name"`
	// ShelterAvailable or not
	ShelterAvailable bool `json:"shelter_available"`
	// Distance to Shelter
	DistanceToShelter float64 `json:"distance_to_shelter"`
}
