package postgres

import (
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/user/config"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(cfg *config.Postgres) *Repository {

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable client_encoding=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.Encoding,
	)

	logLevel := getLogLevel(cfg.LogLevel)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&postgresUser{}, &postgresPlayer{})
	// db.Model(&postgresPlayer{}).Association("user_id")

	return &Repository{
		db: db,
	}
}

// Repository Types

type postgresUser struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// Name of the user
	Name string
	// Email
	Email string `gorm:"unique;not null"`
	// Date of Birth of the player TODO: Fix type
	DateOfBirth string
}

type postgresPlayer struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key"`
	// User is the root entity of player
	UserID uuid.UUID `gorm:"type:uuid;unique not null"`
	// Weight of the player
	Weight float64 `gorm:"<-"`
	// Height of the player
	Height float64 `gorm:"<-"`
	// Preference of the player
	Preference string `gorm:"<-"`
	// GeographicalZone is a group of trails in a region
	ZoneID uuid.UUID
	// CreatedAt is the time when the player registered
	CreatedAt time.Time
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time
}

// Helper for converting to/from domain model

func toAggregate(pu *postgresUser, pp *postgresPlayer) *domain.Player {
	return &domain.Player{
		ID: pp.ID,
		User: &domain.User{
			ID:          pu.ID,
			Email:       pu.Email,
			Name:        pu.Name,
			DateOfBirth: pu.DateOfBirth,
		},
		Weight:     pp.Weight,
		Height:     pp.Height,
		Preference: domain.Preference(pp.Preference),
		ZoneID:     pp.ZoneID,
		CreatedAt:  pp.CreatedAt,
		UpdatedAt:  pp.UpdatedAt,
	}
}

// May be needed later
// func fromAggregate(player *domain.Player) (*postgresUser, *postgresPlayer) {
// 	pu := &postgresUser{
// 		ID:          player.User.ID,
// 		Name:        player.User.Name,
// 		Email:       player.User.Email,
// 		DateOfBirth: player.User.DateOfBirth,
// 	}

// 	pp := &postgresPlayer{
// 		ID:        player.ID,
// 		UserID:    player.User.ID,
// 		Weight:    player.Weight,
// 		Height:    player.Height,
// 		ZoneID:    player.ZoneID,
// 		CreatedAt: player.CreatedAt,
// 		UpdatedAt: player.CreatedAt,
// 	}

// 	return pu, pp
// }

// getLogLevel returns the GORM Log Level
func getLogLevel(l string) logger.LogLevel {
	switch l {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	default:
		return logger.Warn
	}
}
