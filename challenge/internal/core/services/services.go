package services

import (
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/ports"
)

type ChallengeService struct {
	monitor sync.Map
	repo    ports.ChallengeRepository
}

// Factory for creating a new ChallengeService
func NewChallengeService(repo ports.ChallengeRepository) *ChallengeService {
	svc := &ChallengeService{
		repo: repo,
	}

	go svc.MonitorChallenges()

	return svc
}
