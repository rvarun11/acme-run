package workoutsession

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
)

var (
	ErrorPlayerNotFound     = errors.New("the player session not found in repository")
	ErrorFailedToAddPlayer  = errors.New("failed to add the player")
	ErrorUpdatePlayerFailed = errors.New("failed to update player")
)

type PlayerRepository interface {
	Get(uuid.UUID) (aggregate.Player, error)
	Add(aggregate.Player) error
	Update(aggregate.Player) error
}
