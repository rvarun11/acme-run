package postgres

import (
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresChallenge struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"not null"`
	BadgeURL    string
	Criteria    string `gorm:"not null"`
	Goal        float64
	Start       time.Time
	End         time.Time
	CreatedAt   time.Time
}

type postgresBadge struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PlayerID    uuid.UUID `gorm:"not null"`
	ChallengeID uuid.UUID `gorm:"not null"`
	Score       float64
	CompletedOn time.Time
}

type postgresChallengeStats struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	PlayerID        uuid.UUID `gorm:"not null"`
	ChallengeID     uuid.UUID `gorm:"not null"`
	DistanceCovered float64
	EnemiesFought   uint8
	EnemiesEscaped  uint8
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(cfg *config.Postgres) *Repository {

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable client_encoding=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DB_Name,
		cfg.Password,
		cfg.Encoding,
	)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&postgresChallenge{}, &postgresBadge{}, &postgresChallengeStats{})

	return &Repository{
		db: db,
	}
}

// Helper for converting to domain Challenge
func (pc *postgresChallenge) toAggregate() *domain.Challenge {
	return &domain.Challenge{
		ID:          pc.ID,
		Name:        pc.Name,
		Description: pc.Description,
		Criteria:    domain.Criteria(pc.Criteria),
		Goal:        pc.Goal,
		Start:       pc.Start,
		End:         pc.End,
		BadgeURL:    pc.BadgeURL,
		CreatedAt:   pc.CreatedAt,
	}
}

// Helper function to convert to domain Badge
func (pb *postgresBadge) toAggregate(ch *domain.Challenge) *domain.Badge {
	return &domain.Badge{
		ID:          pb.ID,
		PlayerID:    pb.PlayerID,
		Challenge:   ch,
		CompletedOn: pb.CompletedOn,
	}
}

// Helper function to convert to domain Badge
func (pcs *postgresChallengeStats) toAggregate(ch *domain.Challenge) *domain.ChallengeStats {
	return &domain.ChallengeStats{
		PlayerID:        pcs.PlayerID,
		Challenge:       ch,
		DistanceCovered: pcs.DistanceCovered,
		EnemiesFought:   pcs.EnemiesFought,
		EnemiesEscaped:  pcs.EnemiesEscaped,
	}
}
