package services

import (
	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/dto"
	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/ports"
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

func (s *PlayerService) Register(req *dto.PlayerDTO) (*dto.PlayerDTO, error) {
	// TODO: This can be improved because the types of all these fields are same and can cause problems, if ordered incorrectly
	p, err := domain.NewPlayer(req.User.Name, req.User.Email, req.User.DateOfBirth, req.Weight, req.Height)
	if err != nil {
		return &dto.PlayerDTO{}, ports.ErrorCreatePlayerFailed
	}

	player, err := s.repo.Create(p)
	if err != nil {
		return &dto.PlayerDTO{}, ports.ErrorCreatePlayerFailed
	}

	res := dto.FromAggregate(player)
	return res, nil
}

func (s *PlayerService) Get(id uuid.UUID) (*dto.PlayerDTO, error) {
	player, err := s.repo.Get(id)
	if err != nil {
		return &dto.PlayerDTO{}, err
	}

	res := dto.FromAggregate(player)

	return res, nil
}

// func (s *PlayerService) GetByEmail(email string) (*dto.PlayerDTO, error) {
// 	player, err := s.repo.Get(email)
// 	if err != nil {
// 		return &dto.PlayerDTO{}, err
// 	}
// 	playerDTO := dto.ToDTO(player)
// 	return playerDTO, nil
// }

func (s *PlayerService) Update(req *dto.PlayerDTO) (*dto.PlayerDTO, error) {
	player := dto.ToAggregate(req)
	player, err := s.repo.Update(player)
	if err != nil {
		return &dto.PlayerDTO{}, nil
	}
	res := dto.FromAggregate(player)
	return res, nil
}

func (s *PlayerService) List() ([]*dto.PlayerDTO, error) {
	players, err := s.repo.List()
	if err != nil {
		return []*dto.PlayerDTO{}, err
	}

	var playerDTOs []*dto.PlayerDTO
	for _, pp := range players {

		playerDTO := dto.FromAggregate(pp)
		playerDTOs = append(playerDTOs, playerDTO)
	}
	return playerDTOs, nil
}
