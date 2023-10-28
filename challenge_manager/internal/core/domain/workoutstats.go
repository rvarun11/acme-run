package domain

type WorkoutStats struct {
	// Workout Session ID is ID of the workout session. It will allow getting the Player
	WorkoutSessionID string
	// Enemies fought
	EnemiesFought uint8
	// Enemies Escaped
	EnemiesEscaped uint8
	// Distance covered
	DistanceCovered float64
	// TODO: Maybe add stats related to Heart Rate later?
}
