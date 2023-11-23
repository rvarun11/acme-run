package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	"github.com/google/uuid"
)

var (
	// Challenge
	ErrorListChallengesFailed   = errors.New("failed to list challenges")
	ErrNoActiveChallengePresent = errors.New("no active challenges present")
	ErrorChallengeNotFound      = errors.New("the challenge not found in repository")
	ErrorCreateChallengeFailed  = errors.New("failed to add the challenge")
	ErrorUpdateChallengeFailed  = errors.New("failed to update challenge")
	// Challenge Stats
	ErrorChallengeStatsNotFound = errors.New("unable to find challenge stats")
	// Badge
	ErrBadgeNotFound = errors.New("badge not found")
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
	ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Badge, error)
	DispatchBadges(ch *domain.Challenge)
	SubscribeToActiveChallenges(cs *domain.ChallengeStats) error
	ListChallengeStatsByPlayerID(pid uuid.UUID) error
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
	ListBadges() ([]*domain.Badge, error)
	ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Badge, error)
	// ChallengeStats
	CreateOrUpdateChallengeStats(cs *domain.ChallengeStats) error
	ListChallengeStatsByPlayerID(pid uuid.UUID) ([]*domain.ChallengeStats, error)
	ListChallengeStatsByChallengeID(cid uuid.UUID) ([]*domain.ChallengeStats, error)
	// ListEligibleChallengeStatsForChallenge(ch *domain.Challenge) ([]*domain.ChallengeStats, error)
	DeleteChallengeStats(pid uuid.UUID, cid uuid.UUID) error
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
