package services

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Badge Services

// CreateBadge takes a total challenge stats for a player and creates a badge if criteria is met
func (svc *ChallengeService) CreateBadge(cs *domain.ChallengeStats) (*domain.Badge, error) {
	b, err := domain.NewBadge(cs)
	if err != nil {
		return &domain.Badge{}, err
	}
	badge, err := svc.repo.CreateBadge(b)
	if err != nil {
		return &domain.Badge{}, err
	}
	logger.Debug("badge created for challenge stat", zap.Any("stat", cs))

	return badge, nil
}

func (svc *ChallengeService) ListBadges() ([]*domain.Badge, error) {
	badges, err := svc.repo.ListBadges()
	if err != nil {
		return []*domain.Badge{}, err
	}

	return badges, nil
}

func (svc *ChallengeService) ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Badge, error) {
	badges, err := svc.repo.ListBadgesByPlayerID(pid)
	if err != nil {
		return []*domain.Badge{}, err
	}

	return badges, nil
}

func (svc *ChallengeService) CreateOrUpdateChallengeStats(pid uuid.UUID, dc float64, ef uint8, ee uint8, workoutEnd time.Time) error {
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

// ActeFinal runs when a challenge ends and creates badges for all the players who met the critera for the challenge
func (svc *ChallengeService) AssignBadges(ch *domain.Challenge) {
	// 1. Fetch Player Challenge Stats
	csArr, err := svc.repo.ListChallengeStatsByChallengeID(ch.ID)
	if err != nil {
		logger.Fatal("error occured while fetching challenge stats for players of a challenge", zap.Any("challenge", ch))
	}

	// 2. Create Badges, if critera is met
	// var badges []*domain.Badge
	for _, cs := range csArr {
		_, err := svc.CreateBadge(cs)
		logger.Debug("attempting to create badge for stat", zap.Any("stat", cs))
		if err != nil {
			logger.Debug("unable to create badge for challenge stat", zap.Error(err))
			continue
		}
	}

	// 3. Delete challenges stats for the challenge as it has ended

	logger.Info("challenge ended, badges have been assigned", zap.String("challenge", ch.Name))
}
