package ports

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
	"github.com/google/uuid"
)

type ChallengeService interface {
	Create(c *domain.Challenge) (*domain.Challenge, error)
	Get(uuid uuid.UUID) (*domain.Challenge, error)
	Update(player *domain.Challenge) (*domain.Challenge, error)
	List() ([]*domain.Challenge, error)
	ListActive() ([]*domain.Challenge, error)
}

type ChallengeRepository interface {
	Create(c *domain.Challenge) (*domain.Challenge, error)
	Get(uuid uuid.UUID) (*domain.Challenge, error)
	Update(player *domain.Challenge) (*domain.Challenge, error)
	List() ([]*domain.Challenge, error)
	// ListActive() ([]*domain.Challenge, error)
}

// TODO: To be named properly
type BadgeService interface {
	// This service will take in a workout stat and compare it with all active challenges to see if it's meets criteria,
	CheckEligibility(ws *domain.WorkoutStats) (bool, error)
	GetAll(playerID uuid.UUID) ([]*domain.Challenge, error)
}

type BadgeRepository interface {
	// This service will take in

}
