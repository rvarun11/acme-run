package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrorListWorkoutsFailed  = errors.New("failed to list workout")
	ErrorWorkoutNotFound     = errors.New("workout not found in repository")
	ErrorCreateWorkoutFailed = errors.New("failed to create the workout")
	ErrorUpdateWorkoutFailed = errors.New("failed to update workout")
)

type TrailManagerService interface {
	// TODO: List should only return workouts for a particular Player
	// List() ([]*domain.Workout, error)
	// GetTrailManager(workoutID uuid.UUID) (*domain.Workout, error)

	// StartWorkout(workout domain.Workout) error
	// StopWorkout(workout domain.Workout) (*domain.Workout, error)

	// GetWorkoutOptions(workoutID uuid.UUID) (uint8, error)
	// StartWorkoutOption(workoutID uuid.UUID, workoutType uint8) error
	// StopWorkoutOption(workoutID uuid.UUID) error

	// GetDistanceById(workoutID uuid.UUID) (float64, error)
	// GetDistanceCoveredBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (float64, error)
	// GetEscapesMadeById(workoutID uuid.UUID) (uint16, error)
	// GetEscapesMadeBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
	// GetFightsFoughtById(workoutID uuid.UUID) (uint16, error)
	// GetFightsFoughtBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
	// GetSheltersTakenById(workoutID uuid.UUID) (uint16, error)
	// GetSheltersTakenBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
	sendDistance(workoutId uuid.UUID, shelterId uuid.UUID, distance float64)
	retrieveLocation()
	getCloestShelter()
	getShelterDistance() (float64, error)
	getTrail(id uuid.UUID) (*domain.Trail, error)
	getShelter(id uuid.UUID) (*domain.Shelter, error)
}

type TrailRepository interface {
	CreateTrail(name string, startLat, startLong, endLat, endLong float64, closestShelterID uuid.UUID) (uuid.UUID, error)
	UpdateTrailByID(id uuid.UUID, name string, startLat, startLong, endLat, endLong float64, closestShelterID uuid.UUID) error
	DeleteTrailByID(id uuid.UUID) error
	GetTrailByID(id uuid.UUID) (*domain.Trail, error)
}

type ShelterRepository interface {
	CreateShelter(name string, availability bool, lat, long float64) (uuid.UUID, error)
	UpdateShelterByID(id uuid.UUID, name string, availability bool, lat, long float64) error
	DeleteShelterByID(id uuid.UUID) error
	GetShelterByID(id uuid.UUID) (*domain.Shelter, error)
}
