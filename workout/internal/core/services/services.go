package services

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/ports"
	"github.com/google/uuid"
)

type WorkoutService struct {
	repo ports.WorkoutRepository
}

// Factory for creating a new WorkoutService
func NewWorkoutService(repo ports.WorkoutRepository) *WorkoutService {
	return &WorkoutService{
		repo: repo,
	}
}

func (s *WorkoutService) List() ([]*domain.Workout, error) {
	return s.repo.List()
}

func (s *WorkoutService) Get(id uuid.UUID) (*domain.Workout, error) {
	return s.repo.Get(id)
}

// TODO: Add start workout logic here
func (s *WorkoutService) Start(workout domain.Workout) error {
	// this will create the workout
	// Send request Get HRM
	var temp = s.repo.Create(workout)

	return temp
}

func (s *WorkoutService) Pause(id uuid.UUID) (*domain.Workout, error) {
	// TODO: Think of the logic? Should there be a countdown?
	return s.repo.Get(id)
}

func (s *WorkoutService) Stop(id uuid.UUID) (*domain.Workout, error) {
	// Call Update() to update InProgress to False & EndedAt to time.Now()
	var tempWorkout *domain.Workout
	var err error
	tempWorkout, err = s.repo.Get(id)

	// TODO: Better error handling
	if err != nil {
		return nil, err
	}
	// TODO: More logic to find distance covered and other things
	tempWorkout.EndedAt = time.Now()
	tempWorkout.InProgress = false

	s.repo.Update(*tempWorkout)

	return tempWorkout, err
}

func (s *WorkoutService) UpdateHRValue(workoutID uuid.UUID, hrValue uint16) error {
	var tempWorkout *domain.Workout
	var err error
	tempWorkout, err = s.Get(workoutID)

	if err != nil {
		return nil
	}

	tempWorkout.HeartRate = append(tempWorkout.HeartRate, hrValue)
	s.repo.Update(*tempWorkout)
	return nil
}
