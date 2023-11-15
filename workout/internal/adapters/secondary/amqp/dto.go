package amqphandler

import (
	"time"

	"github.com/google/uuid"
)

type challengeStatsDTO struct {
	PlayerID        uuid.UUID `json:"player_id"`
	DistanceCovered float64   `json:"distance_covered"`
	EnemiesFought   uint8     `json:"enemies_fought"`
	EnemiesEscaped  uint8     `json:"enemies_escaped"`
	WorkoutEnd      time.Time `json:"workout_end"`
}
