package http

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	"github.com/google/uuid"
)

type challengeDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BadgeURL    string    `json:"badge_url"`
	Criteria    string    `json:"criteria"`
	Goal        float64   `json:"goal"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	CreatedAt   time.Time `json:"created_at"`
}

// Add ToAggregate
func toAggregate(chDTO *challengeDTO) *domain.Challenge {
	return &domain.Challenge{
		ID:          chDTO.ID,
		Name:        chDTO.Name,
		Description: chDTO.Description,
		Criteria:    domain.Criteria(chDTO.Criteria),
		Goal:        chDTO.Goal,
		Start:       chDTO.Start,
		End:         chDTO.End,
		BadgeURL:    chDTO.BadgeURL,
		CreatedAt:   chDTO.CreatedAt,
	}
}

func fromAggregate(ch *domain.Challenge) *challengeDTO {
	return &challengeDTO{
		ID:          ch.ID,
		Name:        ch.Name,
		Description: ch.Description,
		Criteria:    string(ch.Criteria),
		Goal:        ch.Goal,
		Start:       ch.Start,
		End:         ch.End,
		BadgeURL:    ch.BadgeURL,
		CreatedAt:   ch.CreatedAt,
	}
}
