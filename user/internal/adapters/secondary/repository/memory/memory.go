package memory

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"

	"github.com/google/uuid"
)

type Repository struct {
	players map[uuid.UUID]domain.Player
	sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		players: make(map[uuid.UUID]domain.Player),
	}
}

func (r *Repository) List() ([]*domain.Player, error) {
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

func (r *Repository) Create(player domain.Player) error {
	if r.players == nil {
		r.Lock()
		r.players = make(map[uuid.UUID]domain.Player)
		r.Unlock()
	}

	if _, ok := r.players[player.ID]; ok {
		return fmt.Errorf("player already exist: %w", ports.ErrorCreatePlayerFailed)
	}
	r.Lock()
	r.players[player.ID] = player
	r.Unlock()
	return nil
}

func (mr *Repository) Get(pid uuid.UUID) (*domain.Player, error) {
	if player, ok := mr.players[pid]; ok {
		return &player, nil
	}
	return &domain.Player{}, ports.ErrorPlayerNotFound
}

func (r *Repository) Update(player domain.Player) error {
	if _, ok := r.players[player.ID]; !ok {
		return fmt.Errorf("player does not exist: %w", ports.ErrorUpdatePlayerFailed)
	}
	r.Lock()
	r.players[player.ID] = player
	r.Unlock()
	return nil
}
