package http

import (
	"time"

	"github.com/google/uuid"
)

type ShelterAvailable struct {
	// WorkoutID for which the Shelter Availability is there or not
	WorkoutID uuid.UUID `json:"workout_id"`
	// ShelterAvailable or not
	ShelterID        uuid.UUID `json:"shelter_id"`
	ShelterAvailable bool      `json:"shelter_available"`
	// Distance to Shelter
	DistanceToShelter float64   `json:"distance_to_shelter"`
	ShelterCheckTime  time.Time `json:"shelter_check_time"`
}

type LocationDTO struct {
	// WorkoutID for which the Shelter Availability is there or not
	WorkoutID uuid.UUID `json:"workout_id"`
	// Latitude of the Player
	Latitude float64 `json:"latitude"`
	// Longitude of the Player
	Longitude float64 `json:"longitude"`
	// Time of location
	TimeOfLocation time.Time `json:"time_of_location"`
}

type UserData struct {
	WorkoutID uuid.UUID `json:"workout_id"`
	TrailID   uuid.UUID `json:"trail_id"`
}

type ZoneDTO struct {
	ZoneID   uuid.UUID `json:"zone_id"`
	ZoneName string    `json:"zone_name"`
}

type ShelterDTO struct {
	ShelterID           uuid.UUID `json:"shelter_id"`
	ShelterName         string    `json:"shelter_name"`
	TrailID             uuid.UUID `json:"trail_id"`
	ShelterAvailability bool      `json:"shelter_availability"`
	Longitude           float64   `json:"longitude"`
	Latitude            float64   `json:"latitude"`
}

type TrailDTO struct {
	TrailID        uuid.UUID `json:"trail_id"`
	TrailName      string    `json:"trail_name"`
	ZoneID         uuid.UUID `json:"zone_id"`
	StartLongitude float64   `json:"start_longitude"`
	StartLatitude  float64   `json:"start_latitude"`
	EndLongitude   float64   `json:"end_longitude"`
	EndLatitude    float64   `json:"end_latitude"`
}
