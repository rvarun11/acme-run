package services

import (
	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/domain"

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
	return s.repo.Create(workout)
}

func (s *WorkoutService) Pause(id uuid.UUID) (*domain.Workout, error) {
	// TODO: Think of the logic? Should there be a countdown?
	return s.repo.Get(id)
}

func (s *WorkoutService) Stop(id uuid.UUID) (*domain.Workout, error) {
	// Call Update() to update InProgress to False & EndedAt to time.Now()
	return s.repo.Get(id)
}
