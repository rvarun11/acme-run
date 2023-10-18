package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"

	"github.com/google/uuid"
)

var (
	ErrorListWorkoutsFailed  = errors.New("failed to list workout")
	ErrorWorkoutNotFound     = errors.New("workout not found in repository")
	ErrorCreateWorkoutFailed = errors.New("failed to create the workout")
	ErrorUpdateWorkoutFailed = errors.New("failed to update workout")
)

type WorkoutService interface {
	// TODO: List should only return workouts for a particular Player
	List() ([]*domain.Workout, error)
	Get(id uuid.UUID) (*domain.Workout, error)
	Start(workout domain.Workout) error
	Stop(workout domain.Workout) (*domain.Workout, error)
}

type WorkoutRepository interface {
	List() ([]*domain.Workout, error)
	Create(workout domain.Workout) error
	Get(workout uuid.UUID) (*domain.Workout, error)
	Update(workout domain.Workout) error
}
