package workoutsession

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
)

var (
	ErrWorkoutSessionNotFound       = errors.New("the workout session not found in repository")
	ErrAddWorkoutSessionFailed      = errors.New("failed to add the workout session")
	ErrorUpdateWorkoutSessionFailed = errors.New("failed to update workout session")
)

type WorkoutSessionRepository interface {
	Get(uuid.UUID) (aggregate.WorkoutSession, error)
	Add(aggregate.WorkoutSession) error
	Update(aggregate.WorkoutSession) error
}
