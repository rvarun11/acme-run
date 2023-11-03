package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrorListChallengesFailed  = errors.New("failed to list challenges")
	ErrorChallengeNotFound     = errors.New("the challenge not found in repository")
	ErrorCreateChallengeFailed = errors.New("failed to add the challenge")
	ErrorUpdateChallengeFailed = errors.New("failed to update challenge")
)

type ChallengeService interface {
	Create(ch *domain.Challenge) (*domain.Challenge, error)
	GetByID(cid uuid.UUID) (*domain.Challenge, error)
	Update(ch *domain.Challenge) (*domain.Challenge, error)
	List() ([]*domain.Challenge, error)
}

type ChallengeRepository interface {
	Create(ch *domain.Challenge) (*domain.Challenge, error)
	GetByID(cid uuid.UUID) (*domain.Challenge, error)
	Update(ch *domain.Challenge) (*domain.Challenge, error)
	List() ([]*domain.Challenge, error)
}

// // // TODO: To be named properly
// type BadgeService interface {
// 	// This service will take in a workout stat and compare it with all active challenges to see if it's meets criteria,
// 	Add(cid uuid.UUID, pid uuid.UUID) error
// 	List(playerID uuid.UUID) ([]*domain.Challenge, error)
// }

// type BadgeRepository interface {
// 	// This service will take in
// 	Create(cid uuid.UUID, pid uuid.UUID) error
// 	List(pid uuid.UUID) (*[]domain.Challenge, error)
// }
