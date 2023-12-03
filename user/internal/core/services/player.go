package services

import (
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/ports"
	logger "github.com/CAS735-F23/macrun-teamvsl/user/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	players, err := s.repo.List()
	if err != nil {
		return []*domain.Player{}, err
	}

	return players, nil
}

func (s *PlayerService) Register(req *domain.Player) (*domain.Player, error) {
	// TODO: This can be improved because the types of all these fields are same and can cause problems, if ordered incorrectly
	p, err := domain.NewPlayer(req.User.Name, req.User.Email, req.User.DateOfBirth, req.Weight, req.Height, domain.Preference(req.Preference), req.ZoneID)
	if err != nil {
		return &domain.Player{}, ports.ErrorCreatePlayerFailed
	}

	player, err := s.repo.Create(p)
	if err != nil {
		return &domain.Player{}, ports.ErrorCreatePlayerFailed
	}

	logger.Info("new player registered successfully", zap.String("player", player.User.Name))
	return player, nil
}

func (s *PlayerService) Get(id uuid.UUID) (*domain.Player, error) {
	player, err := s.repo.Get(id)
	if err != nil {
		return &domain.Player{}, err
	}
	logger.Info("retrieved player successfully", zap.String("name", player.User.Name))
	return player, nil
}

// func (s *PlayerService) GetByEmail(email string) (*dto.PlayerDTO, error) {
// 	player, err := s.repo.Get(email)
// 	if err != nil {
// 		return &dto.PlayerDTO{}, err
// 	}
// 	playerDTO := dto.ToDTO(player)
// 	return playerDTO, nil
// }

func (s *PlayerService) Update(req *domain.Player) (*domain.Player, error) {
	player, err := s.repo.Update(req)
	if err != nil {
		return &domain.Player{}, err
	}

	logger.Info("updated player successfully", zap.String("name", player.User.Name))
	return player, nil
}

func (svc *PlayerService) DeletePlayerByID(id uuid.UUID) error {
	// TODO: Can be implemented, if needed
	return nil
}
