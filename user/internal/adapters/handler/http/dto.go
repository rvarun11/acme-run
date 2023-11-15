package http

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"
	"github.com/google/uuid"
)

type userDTO struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID `json:"id"`
	// Name of the user
	Name string `json:"name"`
	// Email
	Email string `json:"email"`
	// DoB
	DateOfBirth string `json:"dob"`
}

type playerDTO struct {
	// ID of the player
	ID uuid.UUID `json:"id"`
	// User is the root entity of player
	User userDTO `json:"user"`
	// Weight of the player
	Weight float64 `json:"weight"`
	// Height of the player
	Height float64 `json:"height"`
	// Preference of the player
	Preference string `json:"preference"`
	// GeographicalZone is a group of trails in a region
	ZoneID string `json:"zone_id"`
	// CreatedAt is the time when the player registered
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time `json:"updated_at"`
}

func (playerDTO *playerDTO) toAggregate() *domain.Player {
	userDTO := domain.User{
		ID:          playerDTO.User.ID,
		Name:        playerDTO.User.Name,
		Email:       playerDTO.User.Email,
		DateOfBirth: playerDTO.User.DateOfBirth,
	}
	return &domain.Player{
		ID:         playerDTO.ID,
		User:       &userDTO,
		Weight:     playerDTO.Weight,
		Height:     playerDTO.Height,
		Preference: domain.Preference(playerDTO.Preference),
		ZoneID:     uuid.MustParse(playerDTO.ZoneID),
		CreatedAt:  playerDTO.CreatedAt,
		UpdatedAt:  playerDTO.UpdatedAt,
	}
}

func fromAggregate(player *domain.Player) *playerDTO {
	return &playerDTO{
		ID: player.ID,
		User: userDTO{
			ID:          player.User.ID,
			Email:       player.User.Email,
			Name:        player.User.Name,
			DateOfBirth: player.User.DateOfBirth,
		},
		Weight:     player.Weight,
		Height:     player.Height,
		Preference: string(player.Preference),
		ZoneID:     player.ZoneID.String(),
		CreatedAt:  player.CreatedAt,
		UpdatedAt:  player.UpdatedAt,
	}
}
