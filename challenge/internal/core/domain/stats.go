package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidChallengeStatEndTime = errors.New("workout end time exceeds challenge end time")

type ChallengeStats struct {
	ID              uuid.UUID
	PlayerID        uuid.UUID
	Challenge       *Challenge
	DistanceCovered float64
	EnemiesFought   uint8
	EnemiesEscaped  uint8
	WorkoutEnd      time.Time
}

func NewChallengeStats(ch *Challenge, pid uuid.UUID, dc float64, ef uint8, ee uint8, workoutEnd time.Time) (*ChallengeStats, error) {
	err := validateTime(ch, workoutEnd)
	if err != nil {
		return &ChallengeStats{}, err
	}
	return &ChallengeStats{
		ID:              uuid.New(),
		PlayerID:        pid,
		Challenge:       ch,
		DistanceCovered: dc,
		EnemiesFought:   ef,
		EnemiesEscaped:  ee,
		WorkoutEnd:      workoutEnd,
	}, nil
}

func (cs *ChallengeStats) GetValidatedScore() (float64, error) {
	switch cs.Challenge.Criteria {
	case DistanceCovered:
		if cs.DistanceCovered >= cs.Challenge.Goal {
			return cs.DistanceCovered, nil
		}
		return 0.0, fmt.Errorf("unable to validate score, got distanced covered=%f for goal=%f", cs.DistanceCovered, cs.Challenge.Goal)
	case FightEnemy:
		if cs.EnemiesFought >= uint8(cs.Challenge.Goal) {
			return float64(cs.EnemiesFought), nil
		}
		return 0.0, fmt.Errorf("unable to validate score, got enemies fought=%d for goal=%f", cs.EnemiesFought, cs.Challenge.Goal)
	case EscapeEnemy:
		if cs.EnemiesEscaped >= uint8(cs.Challenge.Goal) {
			return float64(cs.EnemiesEscaped), nil
		}
		return 0.0, fmt.Errorf("unable to validate score, got enemies escaped=%d for goal=%f", cs.EnemiesEscaped, cs.Challenge.Goal)
	case FightMoreThanEscape:
		if cs.EnemiesFought > cs.EnemiesEscaped {
			return float64(cs.EnemiesFought - cs.EnemiesEscaped), nil
		}
		return 0.0, fmt.Errorf("unable to validate score, got enemies fought=%d and enemies escaped=%d", cs.EnemiesFought, cs.EnemiesEscaped)
	case EscapeMoreThanFight:
		if cs.EnemiesEscaped > cs.EnemiesFought {
			return float64(cs.EnemiesEscaped - cs.EnemiesFought), nil
		}
		return 0.0, fmt.Errorf("unable to validate score, got enemies escaped=%d and enemies fought=%d", cs.EnemiesEscaped, cs.EnemiesFought)
	default:
		return 0.0, fmt.Errorf("cannot validate score for invalid criteria")
	}
}

func validateTime(ch *Challenge, end time.Time) error {
	if end.Before(ch.End) {
		return ErrInvalidChallengeStatEndTime
	}
	return nil
}

// Specific Stats may be needed later

// type DistanceCoveredStats struct {
// 	// Workout Session ID is ID of the workout session. It will allow getting the Player
// 	WorkoutSessionID string
// 	// Distance covered
// 	DistanceCoveredfloat64
// }

// type EnemiesFoughtStats struct {
// 	// Workout Session ID is ID of the workout session. It will allow getting the Player
// 	WorkoutSessionID string
// 	// Enemies fought
// 	EnemiesFought uint8
// }

// type EnemiesEscapedStats struct {
// 	// Workout Session ID is ID of the workout session. It will allow getting the Player
// 	WorkoutSessionID string
// 	// Enemies escaped
// 	EnemiesEscaped uint8
// }

// type EscapeVsFought struct {
// 	WorkoutSessionID string
// 	EnemiesFought    uint8
// 	EnemiesEscaped   uint8
// }

// func (stats *DistanceCoveredStats) Compare(ch *Challenge) (bool, error) {

// 	if ch.Criteria == DistanceCovered {
// 		return stats.DistanceCovered >= ch.Goal, nil
// 	}

// 	return false, fmt.Errorf("unsupported challenge criteria, required DistanceCoveredStats, got %s", ch.Criteria)
// }

// func (stats *EnemiesFoughtStats) Compare(ch *Challenge) (bool, error) {

// 	if ch.Criteria == FightEnemy {
// 		return stats.EnemiesFought >= uint8(ch.Goal), nil
// 	}

// 	return false, fmt.Errorf("unsupported challenge criteria, required EnemiesFoughtStats, got %s", ch.Criteria)
// }

// func (stats *EnemiesEscapedStats) Compare(ch *Challenge) (bool, error) {

// 	if ch.Criteria == EscapeEnemy {
// 		return stats.EnemiesEscaped >= uint8(ch.Goal), nil
// 	}

// 	return false, fmt.Errorf("unsupported challenge criteria, required EnemiesEscaped, got %s", ch.Criteria)
// }

// func (stats *EscapeVsFought) Compare(ch *Challenge) (bool, error) {
// 	if ch.Criteria == FightMoreThanEscape {
// 		return stats.EnemiesFought >= stats.EnemiesEscaped, nil
// 	} else if ch.Criteria == EscapeMoreThanFight {
// 		return stats.EnemiesEscaped >= stats.EnemiesFought, nil
// 	}

// 	return false, fmt.Errorf("unsupported challenge criteria, required FightMoreThanEscape or EscapeMoreThanFight, got %s", ch.Criteria)
// }
