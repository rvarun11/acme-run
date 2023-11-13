package dto

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/domain"
	"github.com/google/uuid"
)

type PlayerDTO struct {
	// ID of the player
	ID uuid.UUID `json:"id"`
	// User is the root entity of player
	User UserDTO `json:"user"`
	// Weight of the player
	Weight float64 `json:"weight"`
	// Height of the player
	Height float64 `json:"height"`
	// Preference of the player
	Preference string `json:"preference"`
	// GeographicalZone is a group of trails in a region
	ZoneID uuid.UUID `json:"zone_id"`
	// CreatedAt is the time when the player registered
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time `json:"updated_at"`
}

// Add ToAggregate
func ToAggregate(playerDTO *PlayerDTO) *domain.Player {
	userDTO := domain.User{
		ID:          playerDTO.User.ID,
		Name:        playerDTO.User.Name,
		Email:       playerDTO.User.Email,
		DateOfBirth: playerDTO.User.DateOfBirth,
	}
	return &domain.Player{
		ID:        playerDTO.ID,
		User:      &userDTO,
		Weight:    playerDTO.Weight,
		Height:    playerDTO.Height,
		ZoneID:    playerDTO.ZoneID,
		CreatedAt: playerDTO.CreatedAt,
		UpdatedAt: playerDTO.UpdatedAt,
	}
}

func FromAggregate(player *domain.Player) *PlayerDTO {
	return &PlayerDTO{
		ID: player.ID,
		User: UserDTO{
			ID:          player.User.ID,
			Email:       player.User.Email,
			Name:        player.User.Name,
			DateOfBirth: player.User.DateOfBirth,
		},
		Weight:    player.Weight,
		Height:    player.Height,
		ZoneID:    player.ZoneID,
		CreatedAt: player.CreatedAt,
		UpdatedAt: player.UpdatedAt,
	}
}
