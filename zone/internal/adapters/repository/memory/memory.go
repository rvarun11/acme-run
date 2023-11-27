package repository

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/ports"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	ts map[uuid.UUID]domain.ZoneManager
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		ts: make(map[uuid.UUID]domain.ZoneManager),
	}
}

func (r *MemoryRepository) AddZoneManagerIntance(t domain.ZoneManager) error {
	if r.ts == nil {
		r.Lock()
		r.ts = make(map[uuid.UUID]domain.ZoneManager)
		r.Unlock()
	}

	if _, ok := r.ts[t.CurrentWorkoutID]; ok {
		return fmt.Errorf("trail manager already connected: %w", ports.ErrorCreateZoneManagerFailed)
	}
	r.Lock()
	r.ts[t.CurrentWorkoutID] = t
	r.Unlock()
	return nil

}

func (r *MemoryRepository) DeleteZoneManagerInstance(wId uuid.UUID) error {
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
		return fmt.Errorf("trail manager with workout ID %v not found: %w", wId, ports.ErrorZoneManagerlNotFound)
	}

	delete(r.ts, keyToDelete)
	return nil
}

func (r *MemoryRepository) Update(t domain.ZoneManager) error {
	if _, ok := r.ts[t.CurrentWorkoutID]; !ok {
		return fmt.Errorf("trail manager does not exist: %w", ports.ErrorUpdateZoneManagerFailed)
	}
	r.Lock()
	r.ts[t.CurrentWorkoutID] = t
	r.Unlock()
	return nil
}

func (r *MemoryRepository) List() ([]*domain.ZoneManager, error) {
	if r.ts == nil {
		// If r.workouts is nil, return an error or handle the case accordingly
		return nil, fmt.Errorf("ps map doesn't exit %w", ports.ErrorListZoneManagerFailed)
	}
	ps := make([]*domain.ZoneManager, 0, len(r.ts))
	for _, t := range r.ts {
		ps = append(ps, &t)
	}
	return ps, nil
}

func (r *MemoryRepository) GetByWorkoutId(wId uuid.UUID) (*domain.ZoneManager, error) {
	for _, t := range r.ts {
		if t.CurrentWorkoutID == wId {
			return &t, nil // Found the peripheral with the matching WorkoutId
		}
	}
	return nil, ports.ErrorZoneManagerlNotFound // No peripheral found with the given WorkoutId
}
