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
		if len(activeChs) == 0 {
			return []*domain.Challenge{}, ports.ErrNoActiveChallengePresent
		}
		return activeChs, nil
	}

	return chs, nil
}

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

	return badge, nil
}

func (svc *ChallengeService) ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Badge, error) {
	badges, err := svc.repo.ListBadgesByPlayerID(pid)
	if err != nil {
		return []*domain.Badge{}, err
	}

	return badges, nil
}

func (svc *ChallengeService) SubscribeToActiveChallenges(pid uuid.UUID, dc float64, ef uint8, ee uint8, workoutEnd time.Time) error {
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

// ActeFinal runs when a challenge ends and creates badges for all the players who met the critera for the challenge
func (svc *ChallengeService) ActeFinal(ch *domain.Challenge) ([]*domain.Badge, error) {
	// 1. Fetch Player Challenge Stats
	csArr, err := svc.repo.ListChallengeStatsByChallengeID(ch.ID)
	if err != nil {
		logger.Fatal("error occured while fetching challenge stats for player of a challenge", zap.Any("challenge", ch))
		return []*domain.Badge{}, err
	}

	// 2. Create Badges, if critera is met
	var badges []*domain.Badge
	for _, cs := range csArr {
		badge, err := svc.CreateBadge(cs)
		if err != nil {
			logger.Debug("unable to create badge for challenge stat", zap.Error(err))
			continue
		}
		badges = append(badges, badge)
	}

	// 3. (optional) Delete challenges stats for the challenge as it has ended

	return badges, nil
}

// This function is runs to monitor active challenges
func (svc *ChallengeService) MonitorChallenges() {
	// Check every 10seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		activeChs, _ := svc.ListChallenges("active")

		for _, ch := range activeChs {
			logger.Debug("found an active challenge, starting monitor")
			go func(ch *domain.Challenge) {
				// Calculate the duration until the end time
				duration := time.Until(ch.End)

				// Sleep until the challenge end time
				time.Sleep(duration)

				// Call HelloWorld when the challenge ends
				svc.ActeFinal(ch)
			}(ch)
		}
	}
}
