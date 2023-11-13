package repository

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	ps map[uuid.UUID]domain.Peripheral
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		ps: make(map[uuid.UUID]domain.Peripheral),
	}
}

func (r *MemoryRepository) AddPeripheralIntance(p domain.Peripheral) error {
	if r.ps == nil {
		r.Lock()
		r.ps = make(map[uuid.UUID]domain.Peripheral)
		r.Unlock()
	}

	if _, ok := r.ps[p.GetID()]; ok {
		return fmt.Errorf("peripheral already connected: %w", ports.ErrorCreatePeripheralFailed)
	}
	r.Lock()
	r.ps[p.GetID()] = p
	r.Unlock()
	return nil

}

func (r *MemoryRepository) DeletePeripheralInstance(wId uuid.UUID) error {
	var keyToDelete uuid.UUID
	found := false

	r.Lock()
	defer r.Unlock()

	for key, p := range r.ps {
		if p.WorkoutId == wId {
			keyToDelete = key
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("peripheral with workout ID %v not found: %w", wId, ports.ErrorPeripheralNotFound)
	}

	delete(r.ps, keyToDelete)
	return nil
}

func (r *MemoryRepository) Get(wid uuid.UUID) (*domain.Peripheral, error) {
	for _, p := range r.ps {
		if p.WorkoutId == workoutId {
			return &p, nil // Found the peripheral with the matching WorkoutId
		}
	}
	return nil, ports.ErrorPeripheralNotFound // No peripheral found with the given WorkoutId
}

func (r *MemoryRepository) Update(p domain.Peripheral) error {
	if _, ok := r.ps[p.GetID()]; !ok {
		return fmt.Errorf("peripheral does not exist: %w", ports.ErrorUpdatePeripheralFailed)
	}
	r.Lock()
	r.ps[p.GetID()] = p
	r.Unlock()
	return nil
}

func (r *MemoryRepository) List() ([]*domain.Peripheral, error) {
	if r.ps == nil {
		// If r.workouts is nil, return an error or handle the case accordingly
		return nil, fmt.Errorf("ps map doesn't exit %w", ports.ErrorListpsFailed)
	}
	ps := make([]*domain.Peripheral, 0, len(r.ps))
	for _, p := range r.ps {
		ps = append(ps, &p)
	}
	return ps, nil
}
