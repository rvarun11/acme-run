// Package entities holds all the entities that are shared across all subdomains
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidWorkout = errors.New("no workout_id matched")
)

// Workout is a entity that represents a workout in all Domains
type Workout struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	WorkoutID uuid.UUID `json:"workout_id"`
	// trailId is the id of the trail player is on
	TrailID uuid.UUID `json:"trail_id"`
	// PlayerID of the player starting the workout session
	PlayerID uuid.UUID `json:"player_id"`
	// InProgress tells whether the workout is in progress
	IsCompleted bool `json:"is_completed"`
	// CreatedAt is the time when the workout was started
	CreatedAt time.Time `json:"created_at"`
	// Duration of the workout
	EndedAt time.Time `json:"ended_at"`
	// EndedAt is the time when the workout was ended
	DistanceCovered float64 `json:"distance_covered"`
	// Player Profile can be either 'cardio' or 'strength'
	Profile string `json:"profile"`
	// HardcoreMode is the difficulty level chosen by the player
	HardcoreMode bool `json:"hardcore_mode"`
	// Shelters taken for a given workout
	Shelters uint8 `json:"shelters_taken"`
	// Fights fought in a given workout
	Fights uint8 `json:"fights_fought"`
	// Escapes made in a given workout
	Escapes uint8 `json:"escapes_made"`
}

type WorkoutOptions struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	WorkoutID uuid.UUID `json:"workout_id"`
	// 1,2,3 for valid current Workout Options, negative for no current workoutOption
	WorkoutOptionsAvailable int8 `json:"workout_options"`
	// 1,2,3 for valid current Workout Options, negative for no current workoutOption
	CurrentWorkoutOption int8 `json:"current_workout_option"`
	// FightsPushDown Ranking bool
	FightsPushDown bool `json:"fights_push_down"`
	// Is Workout Option Active
	IsWorkoutOptionActive bool `json:"is_workout_option_active"`
	// Distance to Shelter
	DistanceToShelter float64 `json:"distance_to_shelter"`
}

func NewWorkout(PlayerID uuid.UUID, TrailID uuid.UUID, HRMID uuid.UUID, HRMConnected bool, hardCoreMode bool) (Workout, error) {
	if PlayerID == uuid.Nil || TrailID == uuid.Nil {
		return Workout{}, ErrInvalidWorkout
	}

	return Workout{
		WorkoutID:       uuid.New(),
		PlayerID:        PlayerID,
		TrailID:         TrailID,
		Profile:         "cardio",
		IsCompleted:     false,
		HardcoreMode:    hardCoreMode,
		CreatedAt:       time.Now(),
		EndedAt:         time.Time{},
		DistanceCovered: 0,
		Shelters:        0,
		Fights:          0,
		Escapes:         0,
	}, nil
}

func (w *Workout) GetID() uuid.UUID {
	return w.WorkoutID
}

func (w *Workout) SetID(id uuid.UUID) {
	w.WorkoutID = id
}

type WorkoutOptionLink struct {
	Name string
	URL  string
}
