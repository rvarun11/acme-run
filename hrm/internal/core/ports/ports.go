package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvs_/player/internal/core/domain"

	"github.com/google/uuid"
)

var (
	ErrorListPlayersFailed  = errors.New("failed to list players")
	ErrorPlayerNotFound     = errors.New("the player session not found in repository")
	ErrorCreatePlayerFailed = errors.New("failed to add the player")
	ErrorUpdatePlayerFailed = errors.New("failed to update player")
)

type PlayerService interface {
	// TODO: I think this place should be able to handle other business functions as well
	List() ([]*domain.Player, error)
	Get(id uuid.UUID) (*domain.Player, error)
	Create(player domain.Player) error
	Update(player domain.Player) (*domain.Player, error)
}

type PlayerRepository interface {
	List() ([]*domain.Player, error)
	Create(player domain.Player) error
	Get(id uuid.UUID) (*domain.Player, error)
	Update(player domain.Player) error
}
