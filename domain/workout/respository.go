package Workout

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
)

var (
	ErrWorkoutNotFound       = errors.New("the workout session not found in repository")
	ErrAddWorkoutFailed      = errors.New("failed to add the workout session")
	ErrorUpdateWorkoutFailed = errors.New("failed to update workout session")
)

type WorkoutRepository interface {
	Get(uuid.UUID) (aggregate.Workout, error)
	Add(aggregate.Workout) error
	Update(aggregate.Workout) error
}
