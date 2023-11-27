package domain

import (
	"time"

	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Badge struct {
	PlayerID    uuid.UUID
	Challenge   *Challenge
	CompletedOn time.Time
	// score is the total score obtained by the player when completing the challenge
	Score float64
}

func NewBadge(cs *ChallengeStats) (*Badge, error) {
	if cs.Challenge.IsActive() {
		logger.Debug("Oops looks like the challenge is not active")
		return &Badge{}, ErrorChallengeIsActive
	}
	score, err := cs.GetValidatedScore()
	if err != nil {
		logger.Debug("Oops looks like the score is not valid")
		return &Badge{}, err
	}
	logger.Debug("The score is", zap.Any("score", score))
	return &Badge{
		Challenge:   cs.Challenge,
		PlayerID:    cs.PlayerID,
		CompletedOn: time.Now(),
		Score:       score,
	}, nil
}
