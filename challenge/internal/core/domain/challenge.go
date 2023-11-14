package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

/*
ChallengeService:
IMPORTANT FOR REPORT:
- For now, the assumption is that there is only one criteria per challenge.

This service is supposed to be called by:
1. By the client when they want to see their badges and on going challenges.
2. By the workout session, so that the challenge manager can send a live notification when they complete a challenge:
	i) When workout sessions starts, it should send the total stats of the player, enemies fought
*/

var (
	ErrorInvalidCriteria   = errors.New("criteria can only be DistanceCovered, Escape, Fight, FightMoreThanEscape or EscapeMoreThanFight")
	ErrorChallengeInactive = errors.New("cannot create a badge as challenge is inactive")
)

type Criteria string

const (
	// Type 1 - Real Time Tracking
	DistanceCovered Criteria = "DistanceCovered"
	EscapeEnemy     Criteria = "Escape"
	FightEnemy      Criteria = "Fight"
	// Type 2 - Can be tracked only when challenge is complete
	FightMoreThanEscape Criteria = "FightMoreThanEscape"
	EscapeMoreThanFight Criteria = "EscapeMoreThanFight"
)

type Challenge struct {
	ID uuid.UUID
	// Name of the challenge, eg. HalloweeK 2023
	Name string
	// Challenge description, eg.
	Description string
	// Badge is the logo received when the challenge is completed
	BadgeURL string
	// Criteria required to complete the challenge
	Criteria Criteria
	// The Goal of the challenge
	Goal float32
	// start time is the time when the challenge starts
	Start time.Time
	// end time is the time when the challenge ends
	End time.Time
	// When the Challenge was Created
	CreatedAt time.Time
}

// NewPlayer is a factory to create a new Player aggregate
func NewChallenge(name string, desc string, badgeUrl string, criteria Criteria, goal float32, start, end time.Time) (*Challenge, error) {
	err := ValidateCriteria(criteria)
	if err != nil {
		return &Challenge{}, err
	}

	challenge := &Challenge{
		ID:          uuid.New(),
		Name:        name,
		Description: desc,
		Criteria:    criteria,
		Goal:        goal,
		Start:       start,
		End:         end,
		BadgeURL:    badgeUrl,
		CreatedAt:   time.Now(),
	}

	return challenge, nil
}

func (ch *Challenge) IsActive() bool {
	currentTime := time.Now()
	return currentTime.After(ch.Start) && currentTime.Before(ch.End)
}

// func GetActiveCriterion()

func ValidateCriteria(c Criteria) error {
	switch c {
	case DistanceCovered, EscapeEnemy, FightEnemy, FightMoreThanEscape, EscapeMoreThanFight:
		return nil
	default:
		return ErrorInvalidCriteria
	}
}
