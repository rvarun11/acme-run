package domain

import (
	"time"

	"github.com/google/uuid"
)

type Badge struct {
	ID          uuid.UUID
	ChallengeID uuid.UUID
	PlayerID    uuid.UUID
	CreatedAt   time.Time
	// score is the total score obtained by the player when completing the challenge
	Score float32
}

func NewBadge(pid uuid.UUID, ch *Challenge, score float32) (*Badge, error) {
	if !ch.IsActive() {
		return &Badge{}, ErrorChallengeInactive
	}
	err := ch.ValidateScore(score)
	if err != nil {
		return &Badge{}, err
	}
	return &Badge{
		ID:          uuid.New(),
		ChallengeID: ch.ID,
		PlayerID:    pid,
		CreatedAt:   time.Now(),
		Score:       score,
	}, nil
}
