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

// Service Interfaces

type ChallengeService interface {
	// Challenges
	CreateChallenge(ch *domain.Challenge) (*domain.Challenge, error)
	GetChallengeByID(cid uuid.UUID) (*domain.Challenge, error)
	UpdateChallenge(ch *domain.Challenge) (*domain.Challenge, error)
	ListChallenges(status string) ([]*domain.Challenge, error)

	// Badges
	CreateBadge(cid uuid.UUID, pid uuid.UUID) error
	ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Challenge, error)

	// StatsTracker
	// Register a player id with the active challenges
	SubscribeToActiveChallenges(cs *domain.ChallengeStats) error

	// RemoveTracker(ch *domain.Challenge) error
	// This notifies all
	// NotifyAll(ch *domain.Challenge) error
}

// type BadgeService interface {
// }

// type StatsTrackerService interface {
// }

// Repository Interfaces

type ChallengeRepository interface {
	// Challenge
	CreateChallenge(ch *domain.Challenge) (*domain.Challenge, error)
	GetChallengeByID(cid uuid.UUID) (*domain.Challenge, error)
	UpdateChallenge(ch *domain.Challenge) (*domain.Challenge, error)
	ListChallenges() ([]*domain.Challenge, error)
	// Badges
	CreateBadge(b *domain.Badge) (*domain.Badge, error)
	ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Challenge, error)
	// ChallengeStats
}

// type ChallengeRepository interface {
// 	CreateChallenge(ch *domain.Challenge) (*domain.Challenge, error)
// 	GetChallengeByID(cid uuid.UUID) (*domain.Challenge, error)
// 	UpdateChallenge(ch *domain.Challenge) (*domain.Challenge, error)
// 	ListChallenges() ([]*domain.Challenge, error)
// }

// type BadgeRepository interface {
// 	CreateBadge(b *domain.Badge) (*domain.Badge, error)
// 	ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Challenge, error)
// }

// type StatsTrackerRepository interface {
// 	// CreateTracker(ws *domain.WorkoutStats) error
// 	// Get
// 	// Delete(ch *domain.Challenge) error
// }
