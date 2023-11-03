package domain

import "fmt"

type DistanceCoveredStats struct {
	// Workout Session ID is ID of the workout session. It will allow getting the Player
	WorkoutSessionID string
	// Distance covered
	DistanceCovered float32
}

type EnemiesFoughtStats struct {
	// Workout Session ID is ID of the workout session. It will allow getting the Player
	WorkoutSessionID string
	// Enemies fought
	EnemiesFought uint8
}

type EnemiesEscapedStats struct {
	// Workout Session ID is ID of the workout session. It will allow getting the Player
	WorkoutSessionID string
	// Enemies escaped
	EnemiesEscaped uint8
}

type EscapeVsFought struct {
	WorkoutSessionID string
	EnemiesFought    uint8
	EnemiesEscaped   uint8
}

func (stats *DistanceCoveredStats) Compare(ch *Challenge) (bool, error) {

	if ch.Criteria == DistanceCovered {
		return stats.DistanceCovered >= ch.Goal, nil
	}

	return false, fmt.Errorf("unsupported challenge criteria, required DistanceCoveredStats, got %s", ch.Criteria)
}

func (stats *EnemiesFoughtStats) Compare(ch *Challenge) (bool, error) {

	if ch.Criteria == FightEnemy {
		return stats.EnemiesFought >= uint8(ch.Goal), nil
	}

	return false, fmt.Errorf("unsupported challenge criteria, required EnemiesFoughtStats, got %s", ch.Criteria)
}

func (stats *EnemiesEscapedStats) Compare(ch *Challenge) (bool, error) {

	if ch.Criteria == EscapeEnemy {
		return stats.EnemiesEscaped >= uint8(ch.Goal), nil
	}

	return false, fmt.Errorf("unsupported challenge criteria, required EnemiesEscaped, got %s", ch.Criteria)
}

func (stats *EscapeVsFought) Compare(ch *Challenge) (bool, error) {
	if ch.Criteria == FightMoreThanEscape {
		return stats.EnemiesFought >= stats.EnemiesEscaped, nil
	} else if ch.Criteria == EscapeMoreThanFight {
		return stats.EnemiesEscaped >= stats.EnemiesFought, nil
	}

	return false, fmt.Errorf("unsupported challenge criteria, required FightMoreThanEscape or EscapeMoreThanFight, got %s", ch.Criteria)
}
