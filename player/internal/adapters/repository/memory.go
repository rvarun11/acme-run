package repository

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvs_/player/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvs_/player/internal/core/domain"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	players map[uuid.UUID]domain.Player
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		players: make(map[uuid.UUID]domain.Player),
	}
}

func (r *MemoryRepository) List() ([]*domain.Player, error) {
	if r.players == nil {
		// If r.players is nil, return an error or handle the case accordingly
		return nil, fmt.Errorf("players map doesn't exit %w", ports.ErrorListPlayersFailed)
	}
	players := make([]*domain.Player, 0, len(r.players))

	for _, player := range r.players {
		players = append(players, &player)
	}
	return players, nil
}

func (r *MemoryRepository) Create(player domain.Player) error {
	if r.players == nil {
		r.Lock()
		r.players = make(map[uuid.UUID]domain.Player)
		r.Unlock()
	}

	if _, ok := r.players[player.GetID()]; ok {
		return fmt.Errorf("player already exist: %w", ports.ErrorCreatePlayerFailed)
	}
	r.Lock()
	r.players[player.GetID()] = player
	r.Unlock()
	return nil
}

func (mr *MemoryRepository) Get(pid uuid.UUID) (*domain.Player, error) {
	if player, ok := mr.players[pid]; ok {
		return &player, nil
	}
	return &domain.Player{}, ports.ErrorPlayerNotFound
}

func (r *MemoryRepository) Update(player domain.Player) error {
	if _, ok := r.players[player.GetID()]; ok {
		return fmt.Errorf("player does not exist: %w", ports.ErrorUpdatePlayerFailed)
	}
	r.Lock()
	r.players[player.GetID()] = player
	r.Unlock()
	return nil
}
