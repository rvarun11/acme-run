// Package entities holds all the entities that are shared across all subdomains
package aggregate

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidWorkoutSession = errors.New("no player associated with the workout session")
)

// Person is a entity that represents a person in all Domains
type WorkoutSession struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	id uuid.UUID
	// PlayerID of the player starting the workout session
	playerID uuid.UUID
	// CreatedAt is the time when the workout was started/created at?
	createdAt time.Time
	// Duration of the workout, TODO: fix type
	duration time.Duration
	// DurationCovered is the total distance covered during the session
	distanceCovered float64
	// isCardio is for determining whether it's a cardio workout or strength workout
	isCardio bool
	// HardcoreMode is the difficulty level chosen by the player
	hardcoreMode bool
	//  HRM Reading from the workout
	heartRate []float64
}

func NewWorkoutSession(playerID uuid.UUID, hardcoreMode bool, isCardio bool) (WorkoutSession, error) {
	if playerID == uuid.Nil {
		return WorkoutSession{}, ErrInvalidWorkoutSession
	}

	return WorkoutSession{
		id:              uuid.New(),
		playerID:        playerID,
		isCardio:        isCardio,
		hardcoreMode:    hardcoreMode,
		createdAt:       time.Now(),
		duration:        0,
		distanceCovered: 0,
		heartRate:       make([]float64, 0),
	}, nil
}

func (ws *WorkoutSession) GetID() uuid.UUID {
	return ws.id
}

func (ws *WorkoutSession) SetID(id uuid.UUID) {
	ws.id = id
}

func (ws *WorkoutSession) GetPlayerID() uuid.UUID {
	return ws.playerID
}

func (ws *WorkoutSession) SetPlayerID(playerID uuid.UUID) {
	ws.playerID = playerID
}

// TODO: Add the rest as per need
