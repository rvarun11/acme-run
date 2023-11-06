package httphandler

import "github.com/google/uuid"

type StartWorkout struct {
	// TrailID chosen by the Player
	TrailID uuid.UUID `json:"trail_id"`
	// PlayerID of the player starting the workout session
	PlayerID uuid.UUID `json:"player_id"`
	// Whether HRM is connected or not
	HRMConnected bool `json:"hrm_connected"`
	// If HRM is connected then HRM ID otherwise garbage
	HRMId uuid.UUID `json:"hrm_id"`
}

type StartWorkoutOption struct {
	// WorkoutID for which the workout option is to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
	// WorkoutID for which the workout option is to be stopped
	Option uint8 `json:"option"`
}

type StopWorkout struct {
	// WorkoutID for the workout to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
}

type StopWorkoutOption struct {
	// WorkoutID for which the workout option is to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
}
