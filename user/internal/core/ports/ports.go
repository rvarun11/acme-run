package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrorListPlayersFailed  = errors.New("failed to list players")
	ErrorPlayerNotFound     = errors.New("the player session not found in repository")
	ErrorCreatePlayerFailed = errors.New("failed to add the player")
	ErrorUpdatePlayerFailed = errors.New("failed to update player")
)

type PlayerService interface {
	Register(player *domain.Player) (*domain.Player, error)
	Get(uuid uuid.UUID) (*domain.Player, error)
	// GetByID(email string) (*domain.Player, error)
	Update(playerDTO *domain.Player) (*domain.Player, error)
	List() ([]*domain.Player, error)
}

type PlayerRepository interface {
	Create(player *domain.Player) (*domain.Player, error)
	Get(uuid uuid.UUID) (*domain.Player, error)
	GetByEmail(email string) (*domain.Player, error)
	Update(player *domain.Player) (*domain.Player, error)
	List() ([]*domain.Player, error)
}
