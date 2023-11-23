package clients

import (
	"time"

	"github.com/google/uuid"
)

type BindPeripheralData struct {
	// PlayerID
	PlayerID uuid.UUID `json:"player_id"`
	// Trail ID
	TrailID uuid.UUID `json:"trail_id"`
	// WorkoutID for the workout to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
	// If HRM is connected then HRM ID otherwise garbage
	HRMId uuid.UUID `json:"hrm_id"`
	// HRM Connected or not
	HRMConnected bool `json:"hrm_connected"`
	// Do we need live location? based on Hardcore mode
	SendLiveLocationToTrailManager bool `json:"send_live_location_to_trail_manager"`
}

type UnbindPeripheralData struct {
	// WorkoutID for the workout to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
}

type AverageHeartRate struct {
	// WorkoutID for the workout to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
	// Average Heart Rate
	AverageHeartRate uint8 `json:"average_heart_rate"`
}

type userDTO struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID `json:"id"`
	// Name of the user
	Name string `json:"name"`
	// Email
	Email string `json:"email"`
	// DoB
	DateOfBirth string `json:"dob"`
}

type playerDTO struct {
	// ID of the player
	ID uuid.UUID `json:"id"`
	// User is the root entity of player
	User userDTO `json:"user"`
	// Weight of the player
	Weight float64 `json:"weight"`
	// Height of the player
	Height float64 `json:"height"`
	// Preference of the player
	Preference string `json:"preference"`
	// GeographicalZone is a group of trails in a region
	ZoneID uuid.UUID `json:"zone_id"`
	// CreatedAt is the time when the player registered
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time `json:"updated_at"`
}
