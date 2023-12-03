package services

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/ports"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
	logger.Info("new challenge created successfully", zap.String("challenge", challenge.Name))

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

func (svc *ChallengeService) DeleteChallengeByID(id uuid.UUID) error {
	// TODO: Can be implemented, if needed
	return nil
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

/*
This is a function to monitor active challenges and create badges
- Note: This should be handled by a cron job. When a challenge ends, the badges should be dispatched.
*/
func (svc *ChallengeService) MonitorChallenges() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		activeChs, _ := svc.ListChallenges("active")

		for _, ch := range activeChs {
			// Check if the challenge is already being monitored
			if _, alreadyMonitoring := svc.monitor.Load(ch.ID); !alreadyMonitoring {
				logger.Debug("found an active challenge, starting monitor")

				// Mark the challenge as being monitored
				svc.monitor.Store(ch.ID, struct{}{})

				// Start a goroutine to monitor the challenge
				go func(ch *domain.Challenge) {
					// Calculate the duration until the end time
					duration := time.Until(ch.End)

					// Sleep until the challenge end time
					time.Sleep(duration)

					// Run ActeFinal
					svc.AssignBadges(ch)

					// Remove the challenge from the monitoring in progress map
					svc.monitor.Delete(ch.ID)
				}(ch)
			}
		}
	}
}
