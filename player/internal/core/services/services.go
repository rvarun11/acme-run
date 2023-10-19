package services

import (
	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/domain"

	"github.com/google/uuid"
)

type PlayerService struct {
	repo ports.PlayerRepository
}

// Factory for creating a new PlayerService
func NewPlayerService(repo ports.PlayerRepository) *PlayerService {
	return &PlayerService{
		repo: repo,
	}
}

func (s *PlayerService) List() ([]*domain.Player, error) {
	return s.repo.List()
}

func (s *PlayerService) Create(player domain.Player) error {
	return s.repo.Create(player)
}

func (s *PlayerService) Get(id uuid.UUID) (*domain.Player, error) {
	return s.repo.Get(id)
}
