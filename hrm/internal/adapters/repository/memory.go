package repository

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/domain"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	hrms map[uuid.UUID]domain.HRM
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		hrms: make(map[uuid.UUID]domain.HRM),
	}
}

func (r *MemoryRepository) List() ([]*domain.Player, error) {
	if r.hrms == nil {
		// If r.players is nil, return an error or handle the case accordingly
		return nil, fmt.Errorf("hrms map doesn't exit %w", ports.ErrorListPlayersFailed)
	}
	hrms := make([]*domain.HRM, 0, len(r.hrms))

	for _, hrm := range r.hrms {
		hrms = append(hrms, &hrm)
	}
	return hrms, nil
}

func (r *MemoryRepository) Create(hrm domain.HRM) error {
	if r.hrms[] == nil {
		r.Lock()
		r.hrms = make(map[uuid.UUID]domain.HRM)
		r.Unlock()
	}

	if _, ok := r.hrms[hrm.GetID()]; ok {
		return fmt.Errorf("player already exist: %w", ports.ErrorCreatePlayerFailed)
	}
	r.Lock()
	r.hrms[hrm.GetID()] = hrm
	r.Unlock()
	return nil
}

func (mr *MemoryRepository) Get(pid uuid.UUID) (*domain.HRM, error) {
	if hrm, ok := mr.hrms[pid]; ok {
		return &hrm, nil
	}
	return &domain.HRM{}, ports.ErrorPlayerNotFound
}

func (r *MemoryRepository) Update(hrm domain.HRM) error {
	if _, ok := r.hrms[hrm.GetID()]; ok {
		return fmt.Errorf("player does not exist: %w", ports.ErrorUpdatePlayerFailed)
	}
	r.Lock()
	r.hrms[hrm.GetID()] = hrm
	r.Unlock()
	return nil
}
