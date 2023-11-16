package repository

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/ports"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	ts map[uuid.UUID]domain.TrailManager
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		ts: make(map[uuid.UUID]domain.TrailManager),
	}
}

func (r *MemoryRepository) AddTrailManagerIntance(t domain.TrailManager) error {
	if r.ts == nil {
		r.Lock()
		r.ts = make(map[uuid.UUID]domain.TrailManager)
		r.Unlock()
	}

	if _, ok := r.ts[t.CurrentWorkoutID]; ok {
		return fmt.Errorf("peripheral already connected: %w", ports.ErrorCreateTrailManagerFailed)
	}
	r.Lock()
	r.ts[t.t.CurrentWorkoutID] = t
	r.Unlock()
	return nil

}

func (r *MemoryRepository) DeleteTrailManagerInstance(wId uuid.UUID) error {
	var keyToDelete uuid.UUID
	found := false

	r.Lock()
	defer r.Unlock()

	for key, t := range r.ts {
		if t.CurrentWorkoutID == wId {
			keyToDelete = key
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("peripheral with workout ID %v not found: %w", wId, ports.ErrorTrailManagerlNotFound)
	}

	delete(r.ts, keyToDelete)
	return nil
}

func (r *MemoryRepository) Update(t domain.TrailManager) error {
	if _, ok := r.ts[t.CurrentWorkoutID]; !ok {
		return fmt.Errorf("peripheral does not exist: %w", ports.ErrorUpdateTrailManagerFailed)
	}
	r.Lock()
	r.ts[t.CurrentWorkoutID] = t
	r.Unlock()
	return nil
}

func (r *MemoryRepository) List() ([]*domain.TrailManager, error) {
	if r.ts == nil {
		// If r.workouts is nil, return an error or handle the case accordingly
		return nil, fmt.Errorf("ps map doesn't exit %w", ports.ErrorListTrailManagerFailed)
	}
	ps := make([]*domain.TrailManager, 0, len(r.ts))
	for _, t := range r.ts {
		ps = append(ps, &t)
	}
	return ps, nil
}

func (r *MemoryRepository) GetByWorkoutId(wId uuid.UUID) (*domain.TrailManager, error) {
	for _, t := range r.ts {
		if t.WorkoutId == wId {
			return &t, nil // Found the peripheral with the matching WorkoutId
		}
	}
	return nil, ports.ErrorTrailManagerlNotFound // No peripheral found with the given WorkoutId
}
