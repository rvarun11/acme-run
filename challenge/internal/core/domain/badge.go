package domain

import (
	"time"

	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
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
		return &Badge{}, ErrorChallengeIsActive
	}
	score, err := cs.GetValidatedScore()
	if err != nil {
		return &Badge{}, err
	}
	logger.Debug("badge score", zap.Any("score", score))
	return &Badge{
		Challenge:   cs.Challenge,
		PlayerID:    cs.PlayerID,
		CompletedOn: time.Now(),
		Score:       score,
	}, nil
}
