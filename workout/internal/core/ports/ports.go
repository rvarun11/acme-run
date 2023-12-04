package ports

import (
	"errors"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"

	"github.com/google/uuid"
)

var (
	ErrorListWorkoutsFailed         = errors.New("failed to list workout")
	ErrorWorkoutNotFound            = errors.New("workout not found in repository")
	ErrorCreateWorkoutFailed        = errors.New("failed to create the workout")
	ErrorActiveWorkoutAlreadyExists = errors.New("active workout already exists")
	ErrorUpdateWorkoutFailed        = errors.New("failed to update workout")
	ErrorWorkoutOptionUnavailable   = errors.New("workout option unavailable")
	ErrorWorkoutOptionInvalid       = errors.New("workout option invalid")
	ErrInvalidWorkout               = errors.New("no workout_id matched")
	ErrWorkoutOptionAlreadyActive   = errors.New("workout option is already active")
	ErrWorkoutOptionAlreadyInActive = errors.New("no workout option is active")
	ErrWorkoutAlreadyCompleted      = errors.New("workout already completed")
)

type WorkoutService interface {
	List() ([]*domain.Workout, error)
	GetWorkout(workoutID uuid.UUID) (*domain.Workout, error)

	StartWorkout(workout domain.Workout) error
	StopWorkout(workout domain.Workout) (*domain.Workout, error)

	GetWorkoutOptions(workoutID uuid.UUID) (uint8, error)
	StartWorkoutOption(workoutID uuid.UUID, option string) (string, error)
	StopWorkoutOption(workoutID uuid.UUID) (string, error)

	UpdateDistanceTravelled(workoutID uuid.UUID, latitude float64, longitude float64, timeOfLocation time.Time) error
	UpdateShelter(workoutID uuid.UUID, DistanceToShelter float64) error
	ComputeWorkoutOptionsOrder() error

	GetDistanceById(workoutID uuid.UUID) (float64, error)
	GetDistanceCoveredBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (float64, error)
	GetEscapesMadeById(workoutID uuid.UUID) (uint16, error)
	GetEscapesMadeBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
	GetFightsFoughtById(workoutID uuid.UUID) (uint16, error)
	GetFightsFoughtBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
	GetSheltersTakenById(workoutID uuid.UUID) (uint16, error)
	GetSheltersTakenBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
}

type WorkoutRepository interface {
	//List() ([]*domain.Workout, error)
	Create(workout *domain.Workout, workoutOptions *domain.WorkoutOptions) error

	GetWorkout(workoutID uuid.UUID) (*domain.Workout, error)
	UpdateWorkout(workout *domain.Workout) (*domain.Workout, error)
	GetWorkoutOptions(workoutID uuid.UUID) (*domain.WorkoutOptions, error)
	UpdateWorkoutOptions(workoutOptions *domain.WorkoutOptions) (*domain.WorkoutOptions, error)

	DeleteWorkoutOptions(workoutID uuid.UUID) error

	GetDistanceByID(workoutID uuid.UUID) (float64, error)
	GetDistanceCoveredBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (float64, error)
	GetEscapesMadeByID(workoutID uuid.UUID) (uint16, error)
	GetEscapesMadeBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
	GetFightsFoughtByID(workoutID uuid.UUID) (uint16, error)
	GetFightsFoughtBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
	GetSheltersTakenByID(workoutID uuid.UUID) (uint16, error)
	GetSheltersTakenBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error)
}

type WorkoutStatsPublisher interface {
	PublishWorkoutStats(workoutStats *domain.Workout) error
}

type UserServiceClient interface {
	GetWorkoutPreferenceOfUser(playerID uuid.UUID) (string, error)
	GetUserAge(playerID uuid.UUID) (uint8, error)
}

type PeripheralClient interface {
	GetAverageHeartRateOfUser(workoutID uuid.UUID) (uint8, error)
	BindPeripheralData(trailID uuid.UUID, playerID uuid.UUID, workoutID uuid.UUID, hrmID uuid.UUID, HRMConnected bool, SendLiveLocationToTrailManager bool) error
	UnbindPeripheralData(workoutID uuid.UUID) error
}
