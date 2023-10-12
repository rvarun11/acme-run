package memory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
	PlayerRepository "github.com/rvarun11/macrun-teamvs/domain/player"
)

type MemoryRepository struct {
	players map[uuid.UUID]aggregate.Player
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		players: make(map[uuid.UUID]aggregate.Player),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (aggregate.Player, error) {
	if player, ok := mr.players[id]; ok {
		return player, nil
	}
	return aggregate.Player{}, PlayerRepository.ErrorPlayerNotFound
}

func (mr *MemoryRepository) Add(player aggregate.Player) error {
	if mr.players == nil {
		mr.Lock()
		mr.players = make(map[uuid.UUID]aggregate.Player)
		mr.Unlock()
	}
	if _, ok := mr.players[player.GetID()]; ok {
		return fmt.Errorf("player already exist: %w", PlayerRepository.ErrorFailedToAddPlayer)
	}
	mr.Lock()
	mr.players[player.GetID()] = player
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(player aggregate.Player) error {
	if _, ok := mr.players[player.GetID()]; ok {
		return fmt.Errorf("workout session does not exist: %w", PlayerRepository.ErrorUpdatePlayerFailed)
	}
	mr.Lock()
	mr.players[player.GetID()] = player
	mr.Unlock()
	return nil
}
