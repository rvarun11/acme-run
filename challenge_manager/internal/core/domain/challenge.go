package domain

import (
	"time"
)

/*
This service is supposed to be called by:
1. By the client when they want to see their badges and on going challenges
2. By the workout session, so that the challenge manager can send a live notification when they complete a challenge:
	i) When workout sessions starts, it should send the total stats of the player, enemies fought

*/

type ChallengeCriteron struct {
	// Enemies fought
	EnemiesFought uint8
	// Enemies escaped
	EnemiesEscaped uint8
	// Distance covered
	DistanceCovered float64
}

type Challenge struct {
	// Name of the challenge, eg. HalloweeK 2023
	Name string
	// Challenge description, eg.
	Description string
	// Criterion required to complete challenge
	Criteria ChallengeCriteron
	// start time is the time when the challenge starts
	StartTime time.Time
	// end time is the time when the challenge ends
	EndTime time.Time
	// Badge is the logo received when the challenge is completed
	BadgeURL string
}

// NewPlayer is a factory to create a new Player aggregate
func NewChallenge(name string, desc string, start time.Time, end time.Time, badgeUrl string) (*Challenge, error) {

	challenge := &Challenge{
		Name:        name,
		Description: desc,
		StartTime:   start,
		EndTime:     end,
		BadgeURL:    badgeUrl,
	}

	return challenge, nil
}
