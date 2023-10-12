package repository

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/domain"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	workouts map[uuid.UUID]domain.Workout
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		workouts: make(map[uuid.UUID]domain.Workout),
	}
}

func (r *MemoryRepository) List() ([]*domain.Workout, error) {
	if r.workouts == nil {
		// If r.workouts is nil, return an error or handle the case accordingly
		return nil, fmt.Errorf("workouts map doesn't exit %w", ports.ErrorListWorkoutsFailed)
	}
	workouts := make([]*domain.Workout, 0, len(r.workouts))

	for _, workout := range r.workouts {
		workouts = append(workouts, &workout)
	}
	return workouts, nil
}

func (r *MemoryRepository) Create(workout domain.Workout) error {
	if r.workouts == nil {
		r.Lock()
		r.workouts = make(map[uuid.UUID]domain.Workout)
		r.Unlock()
	}

	if _, ok := r.workouts[workout.GetID()]; ok {
		return fmt.Errorf("workout already exist: %w", ports.ErrorCreateWorkoutFailed)
	}
	r.Lock()
	r.workouts[workout.GetID()] = workout
	r.Unlock()
	return nil
}

func (mr *MemoryRepository) Get(pid uuid.UUID) (*domain.Workout, error) {
	if workout, ok := mr.workouts[pid]; ok {
		return &workout, nil
	}
	return &domain.Workout{}, ports.ErrorWorkoutNotFound
}

func (r *MemoryRepository) Update(workout domain.Workout) error {
	if _, ok := r.workouts[workout.GetID()]; ok {
		return fmt.Errorf("workout does not exist: %w", ports.ErrorUpdateWorkoutFailed)
	}
	r.Lock()
	r.workouts[workout.GetID()] = workout
	r.Unlock()
	return nil
}
