package amqp

import (
	"github.com/google/uuid"
)

type challengeStatsDTO struct {
	PlayerID        uuid.UUID `json:"player_id"`
	DistanceCovered float32   `json:"distance_covered"`
	EnemiesFought   uint8     `json:"enemies_fought"`
	EnemiesEscaped  uint8     `json:"enemies_escaped"`
}
