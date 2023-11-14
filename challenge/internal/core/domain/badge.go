package domain

import (
	"time"

	"github.com/google/uuid"
)

type Badge struct {
	ID          uuid.UUID
	PlayerID    uuid.UUID
	Challenge   *Challenge
	CompletedOn time.Time
	// score is the total score obtained by the player when completing the challenge
	Score float64
}

func NewBadge(cs *ChallengeStats) (*Badge, error) {
	if !cs.Challenge.IsActive() {
		return &Badge{}, ErrorChallengeInactive
	}
	score, err := cs.GetValidatedScore()
	if err != nil {
		return &Badge{}, err
	}

	return &Badge{
		ID:          uuid.New(),
		Challenge:   cs.Challenge,
		PlayerID:    cs.PlayerID,
		CompletedOn: time.Now(),
		Score:       score,
	}, nil
}
