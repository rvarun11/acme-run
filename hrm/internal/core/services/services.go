package services

import (
	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/domain"
)

type HRMService struct {
	repo ports.HRMRepository
}

// Factory for creating a new PlayerService
func NewHRMService(repo ports.HRMRepository) *HRMService {
	return &HRMService{
		repo: repo,
	}
}

func (s *HRMService) List() ([]*domain.HRM, error) {
	return s.repo.List()
}

func (s *HRMService) Create(hrm domain.HRM) error {
	return s.repo.Create(hrm)
}

func (s *HRMService) Get(id string) (*domain.HRM, error) {
	return s.repo.Get(id)
}
