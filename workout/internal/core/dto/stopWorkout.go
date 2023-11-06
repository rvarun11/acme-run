package dto

import "github.com/google/uuid"

type StopWorkout struct {
	// WorkoutID for the workout to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
}

type StopWorkoutOption struct {
	// WorkoutID for which the workout option is to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
}
