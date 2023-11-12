package clients

import "github.com/google/uuid"

type BindPeripheralData struct {
	// PlayerID
	PlayerID uuid.UUID `json:"player_id"`
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
