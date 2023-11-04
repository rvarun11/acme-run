package services

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/ports"
	"github.com/google/uuid"
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

func (svc *ChallengeService) Create(req *domain.Challenge) (*domain.Challenge, error) {
	ch, err := domain.NewChallenge(req.Name, req.Description, req.BadgeURL, req.Criteria, req.Goal, req.Start, req.End)
	if err != nil {
		return &domain.Challenge{}, ports.ErrorCreateChallengeFailed
	}

	challenge, err := svc.repo.Create(ch)
	if err != nil {
		return &domain.Challenge{}, ports.ErrorCreateChallengeFailed
	}

	return challenge, nil

}

func (svc *ChallengeService) GetByID(id uuid.UUID) (*domain.Challenge, error) {
	ch, err := svc.repo.GetByID(id)
	if err != nil {
		return &domain.Challenge{}, err
	}

	return ch, nil
}

func (svc *ChallengeService) Update(req *domain.Challenge) (*domain.Challenge, error) {
	ch, err := svc.repo.Update(req)
	if err != nil {
		return &domain.Challenge{}, nil
	}

	return ch, nil
}

func (svc *ChallengeService) List(status string) ([]*domain.Challenge, error) {
	chs, err := svc.repo.List()
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

func (svc *ChallengeService) CreateBadge(pid uuid.UUID, ch *domain.Challenge, score float32) (*domain.Badge, error) {
	b, err := domain.NewBadge(pid, ch, score)
	if err != nil {
		return &domain.Badge{}, err
	}
	badge, err := svc.repo.CreateBadge(b)
	if err != nil {
		return &domain.Badge{}, err
	}

	return badge, nil
}

// This should be part of the Badge Service
func (svc *ChallengeService) X(pid uuid.UUID, ch *domain.Challenge, score float32) (*domain.Badge, error) {
	// b, err := domain.NewBadge(pid, ch, score)
	// if err != nil {
	// 	return &domain.Badge{}, err
	// }
	// badge, err := svc.repo.CreateBadge(b)
	// if err != nil {
	// 	return &domain.Badge{}, err
	// }

	return &domain.Badge{}, nil
}
