package services

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/ports"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChallengeService struct {
	repo ports.ChallengeRepository
}

// Factory for creating a new ChallengeService
func NewChallengeService(repo ports.ChallengeRepository) *ChallengeService {
	return &ChallengeService{
		repo: repo,
	}
}

// Challenge Services

func (svc *ChallengeService) CreateChallenge(req *domain.Challenge) (*domain.Challenge, error) {
	ch, err := domain.NewChallenge(req.Name, req.Description, req.BadgeURL, req.Criteria, req.Goal, req.Start, req.End)
	if err != nil {
		return &domain.Challenge{}, ports.ErrorCreateChallengeFailed
	}

	challenge, err := svc.repo.CreateChallenge(ch)
	if err != nil {
		return &domain.Challenge{}, ports.ErrorCreateChallengeFailed
	}
	logger.Info("created new challenge", zap.String("challenge", challenge.Name))

	return challenge, nil

}

func (svc *ChallengeService) GetChallengeByID(id uuid.UUID) (*domain.Challenge, error) {
	ch, err := svc.repo.GetChallengeByID(id)
	if err != nil {
		return &domain.Challenge{}, err
	}

	return ch, nil
}

func (svc *ChallengeService) UpdateChallenge(req *domain.Challenge) (*domain.Challenge, error) {
	ch, err := svc.repo.UpdateChallenge(req)
	if err != nil {
		return &domain.Challenge{}, nil
	}

	return ch, nil
}

func (svc *ChallengeService) ListChallenges(status string) ([]*domain.Challenge, error) {
	chs, err := svc.repo.ListChallenges()
	if err != nil {
		return []*domain.Challenge{}, err
	}

	if status == "active" {
		var activeChs []*domain.Challenge
		for _, ch := range chs {
			if ch.IsActive() {
				activeChs = append(activeChs, ch)
			}
		}
		return activeChs, nil
	}

	return chs, nil
}

// Badge Services

func (svc *ChallengeService) CreateBadge(cs *domain.ChallengeStats) (*domain.Badge, error) {
	b, err := domain.NewBadge(cs)
	if err != nil {
		return &domain.Badge{}, err
	}
	badge, err := svc.repo.CreateBadge(b)
	if err != nil {
		return &domain.Badge{}, err
	}

	return badge, nil
}

// TODO: This should be part of the Badge Service
func (svc *ChallengeService) ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Badge, error) {
	badges, err := svc.repo.ListBadgesByPlayerID(pid)
	if err != nil {
		return []*domain.Badge{}, err
	}

	return badges, nil
}

func (svc *ChallengeService) SubscribeToActiveChallenges(pid uuid.UUID, dc float64, ef uint8, ee uint8, workoutEnd time.Time) error {
	// TODO: This should be in the context
	activeChs, err := svc.ListChallenges("active")
	if err != nil {
		return err
	}

	for _, ch := range activeChs {
		cs, err := domain.NewChallengeStats(ch, pid, dc, ef, ee, workoutEnd)
		if err != nil {
			logger.Debug("cannot create challenge stat", zap.Error(err))
			continue
		}
		err = svc.repo.CreateOrUpdateChallengeStats(cs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *ChallengeService) ListChallengeStatsByPlayerID(pid uuid.UUID) ([]*domain.ChallengeStats, error) {
	csArr, err := svc.repo.ListChallengeStatsByPlayerID(pid)
	if err != nil {
		return []*domain.ChallengeStats{}, err
	}
	return csArr, nil
}

// This function runs when a challenge ends, TODO: Rename once you have more clarity
func (svc *ChallengeService) ActeFinal(ch *domain.Challenge) ([]*domain.Badge, error) {
	// 1. Fetch Challenge Stats: NOTE: This should be further be improved by fetching eligible stats directly at repo level
	csArr, err := svc.repo.ListChallengeStatsByChallengeID(ch.ID)
	if err != nil {
		logger.Fatal("error occured while fetching challenge stats for challenge", zap.String("name", ch.Name))
		return []*domain.Badge{}, err
	}
	// 2. Create Badge
	var badges []*domain.Badge
	for _, cs := range csArr {
		badge, err := domain.NewBadge(cs)
		if err != nil {
			logger.Debug("unable to create badge for challenge stat")
			return []*domain.Badge{}, err
		}
		badges = append(badges, badge)
	}
	// 3. Delete all Challenge Stats
	// TODO: Delete all challenge stats, once the badges are created and the challenge ends.
	return badges, nil
}
