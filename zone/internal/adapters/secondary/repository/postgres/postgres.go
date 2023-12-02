package postgres

import (
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var cfg *config.AppConfiguration = config.Config

type postgresShelter struct {
	ShelterID           uuid.UUID `gorm:"type:uuid;primaryKey;unique"`
	ShelterName         string    `gorm:"type:string;not null;unique"`
	ShelterAvailability bool
	TrailID             uuid.UUID
	Longitude           float64
	Latitude            float64
}
type postgresTrail struct {
	TrailID        uuid.UUID `gorm:"type:uuid;primaryKey;unique"`
	TrailName      string    `gorm:"type:string;not null;unique"`
	ZoneID         uuid.UUID `gorm:"not null"`
	StartLongitude float64
	StartLatitude  float64
	EndLongitude   float64
	EndLatitude    float64
	CreatedAt      time.Time `gorm:"type:timestamp"`
}

type postgresZone struct {
	ZoneID   uuid.UUID `gorm:"type:uuid;primaryKey;unique"`
	ZoneName string    `gorm:"type:string;not null;unique"`
}

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(cfg *config.Postgres) *DBRepository {

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
	db.AutoMigrate(&postgresTrail{}, &postgresShelter{}, &postgresZone{})
	return &DBRepository{db: db}
}

func (ptrail *postgresTrail) toAggregate() *domain.Trail {

	return &domain.Trail{
		TrailID:        ptrail.TrailID,
		TrailName:      ptrail.TrailName,
		ZoneID:         ptrail.ZoneID,
		StartLongitude: ptrail.StartLongitude,
		StartLatitude:  ptrail.StartLatitude,
		EndLongitude:   ptrail.EndLongitude,
		EndLatitude:    ptrail.EndLatitude,
		CreatedAt:      ptrail.CreatedAt,
	}
}
func (pshelter *postgresShelter) toAggregate() *domain.Shelter {

	return &domain.Shelter{
		ShelterID:           pshelter.ShelterID,
		ShelterName:         pshelter.ShelterName,
		TrailID:             pshelter.TrailID,
		ShelterAvailability: pshelter.ShelterAvailability,
		Longitude:           pshelter.Longitude,
		Latitude:            pshelter.Latitude,
	}
}

func (pzone *postgresZone) toAggregate() *domain.Zone {

	return &domain.Zone{
		ZoneID:   pzone.ZoneID,
		ZoneName: pzone.ZoneName,
	}
}
