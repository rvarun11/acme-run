// Package entities holds all the entities that are shared across all subdomains
package aggregate

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidWorkout = errors.New("no player associated with the workout session")
)

// Workout is a entity that represents a workout in all Domains
type Workout struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	id uuid.UUID
	// trailId is the id of the trail player is on
	trailId uuid.UUID
	// PlayerID of the player starting the workout session
	playerID uuid.UUID
	// CreatedAt is the time when the workout was started/created at?
	createdAt time.Time
	// Duration of the workout, TODO: fix type
	endedAt *time.Time
	// DurationCovered is the total distance covered during the session
	distanceCovered float64
	// TODO: temp value. It can be either "cardio", "physical" or "dynamic"
	option string
	// HardcoreMode is the difficulty level chosen by the player
	hardcoreMode bool
	//  HRM Reading from the workout
	// heartRate []valueobject.HeartRate
}

func NewWorkout(playerID uuid.UUID, trailId uuid.UUID, hardcoreMode bool, option string) (Workout, error) {
	if playerID == uuid.Nil {
		return Workout{}, ErrInvalidWorkout
	}

	return Workout{
		id:              uuid.New(),
		playerID:        playerID,
		trailId:         trailId,
		option:          option,
		hardcoreMode:    hardcoreMode,
		createdAt:       time.Now(),
		endedAt:         nil,
		distanceCovered: 0,
		// heartRates:      make([]valueobject.HeartRate, 0),
	}, nil
}

func (ws *Workout) GetID() uuid.UUID {
	return ws.id
}

func (ws *Workout) SetID(id uuid.UUID) {
	ws.id = id
}

func (ws *Workout) GetPlayerID() uuid.UUID {
	return ws.playerID
}

func (ws *Workout) SetPlayerID(playerID uuid.UUID) {
	ws.playerID = playerID
}

// TODO: Add the rest as per need
