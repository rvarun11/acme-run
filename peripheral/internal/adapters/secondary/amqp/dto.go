package rabbitmqhandler

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
	TimeOfLocation time.Time `json:"time_of_reading"`
}

type BindPeripheralData struct {
	// PlayerID
	PlayerID uuid.UUID `json:"player_id"`
	// WorkoutID for the workout to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
	// if HRM Connected is false then no HRM mock
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

type HRMDTO struct {
	PlayerID uuid.UUID `json:"player_id"`
	HRMID    uuid.UUID `json:"hrm_id"`
	Connect  bool      `json:"connect"`
}
