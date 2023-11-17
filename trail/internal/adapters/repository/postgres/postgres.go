package postgres

import (
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/trail/config"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var cfg *config.AppConfiguration = config.Config

type TrailRepository struct {
	db *gorm.DB
}

type ShelterRepository struct {
	db *gorm.DB
}

type ZoneRepository struct {
	db *gorm.DB
}

func NewTrailRepository(cfg *config.Postgres) *TrailRepository {

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
	db.AutoMigrate(&postgresTrail{})
	return &TrailRepository{db: db}
}

func NewZoneRepository(cfg *config.Postgres) *ZoneRepository {

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
	db.AutoMigrate(&postgresTrail{})
	return &ZoneRepository{db: db}
}

func NewShelterRepository(cfg *config.Postgres) *ShelterRepository {

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
	db.AutoMigrate(&postgresShelter{})
	return &ShelterRepository{db: db}
}

// Repository Types

// Structs for GORM
type postgresTrail struct {
	TrailID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TrailName      string    `gorm:"type:string;not null"`
	ZoneID         uuid.UUID `gorm:"not null"`
	StartLongitude float64
	StartLatitude  float64
	EndLongitude   float64
	EndLatitude    float64
	CreatedAt      time.Time `gorm:"type:timestamp"`
}

type postgresShelter struct {
	ShelterID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ShelterName         string    `gorm:"type:string;not null"`
	ShelterAvailability bool
	TrailID             uuid.UUID
	Longitude           float64
	Latitude            float64
}

type postgresZone struct {
	ZoneID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	ZoneName string    `gorm:"type:string;not null"`
}

// Override the TableName method to specify the custom table name for the Trail model
func (postgresTrail) TableName() string {
	return cfg.Postgres.DB_Name + ".trail"
}

// Override the TableName method to specify the custom table name for the Shelter model
func (postgresShelter) TableName() string {
	return cfg.Postgres.DB_Name + ".shelter"
}

// Override the TableName method to specify the custom table name for the Zone model
func (postgresZone) TableName() string {
	return cfg.Postgres.DB_Name + ".zone"
}

func toTrailAggregate(ptrail *postgresTrail) *domain.Trail {

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

func toShelterAggregate(pshelter *postgresShelter) *domain.Shelter {

	return &domain.Shelter{
		ShelterID:           pshelter.ShelterID,
		ShelterName:         pshelter.ShelterName,
		TrailID:             pshelter.TrailID,
		ShelterAvailability: pshelter.ShelterAvailability,
		Longitude:           pshelter.Longitude,
		Latitude:            pshelter.Latitude,
	}
}

func toTrailPostgres(trail *domain.Trail) *postgresTrail {

	return &postgresTrail{
		TrailID:        trail.TrailID,
		TrailName:      trail.TrailName,
		ZoneID:         trail.ZoneID,
		StartLongitude: trail.StartLongitude,
		StartLatitude:  trail.StartLatitude,
		EndLongitude:   trail.EndLongitude,
		EndLatitude:    trail.EndLatitude,
		CreatedAt:      trail.CreatedAt,
	}
}

func toShelterPostgres(shelter *domain.Shelter) *postgresShelter {

	return &postgresShelter{
		ShelterID:           shelter.ShelterID,
		ShelterName:         shelter.ShelterName,
		ShelterAvailability: shelter.ShelterAvailability,
		Longitude:           shelter.Longitude,
		Latitude:            shelter.Latitude,
	}
}

// Repository Functions

func (repo *TrailRepository) CreateTrail(tId uuid.UUID, name string, zId uuid.UUID, startLat, startLong, endLat, endLong float64) (uuid.UUID, error) {
	trail := postgresTrail{
		TrailID:        tId,
		TrailName:      name,
		ZoneID:         zId,
		StartLatitude:  startLat,
		StartLongitude: startLong,
		EndLatitude:    endLat,
		EndLongitude:   endLong,
		CreatedAt:      time.Now(),
	}
	if err := repo.db.Create(&trail).Error; err != nil {
		return uuid.Nil, err
	}
	return tId, nil
}

func (repo *TrailRepository) UpdateTrailByID(id uuid.UUID, name string, zId uuid.UUID, startLat, startLong, endLat, endLong float64) error {
	return repo.db.Model(&postgresTrail{}).Where("trail_id = ?", id).Updates(postgresTrail{
		TrailName:      name,
		ZoneID:         zId,
		StartLatitude:  startLat,
		StartLongitude: startLong,
		EndLatitude:    endLat,
		EndLongitude:   endLong,
	}).Error
}

func (repo *TrailRepository) DeleteTrailByID(id uuid.UUID) error {
	return repo.db.Delete(&postgresTrail{}, "trail_id = ?", id).Error
}

func (repo *TrailRepository) GetTrailByID(id uuid.UUID) (*domain.Trail, error) {
	var trail postgresTrail
	if err := repo.db.Where("trail_id = ?", id).First(&trail).Error; err != nil {
		return nil, err // remove the domain.Trail type
	}
	return toTrailAggregate(&trail), nil
}

func (repo *TrailRepository) List() ([]*domain.Trail, error) {
	var postgresTrails []postgresTrail
	if err := repo.db.Find(&postgresTrails).Error; err != nil {
		return nil, err
	}

	// Convert the postgresTrail records to domain.Trail objects
	domainTrails := make([]*domain.Trail, len(postgresTrails))
	for i, ptrail := range postgresTrails {
		domainTrails[i] = toTrailAggregate(&ptrail)
	}

	return domainTrails, nil
}
func (repo *TrailRepository) ListTrailsByZoneId(zId uuid.UUID) ([]*domain.Trail, error) {
	var postgresTrails []postgresTrail
	if err := repo.db.Where("zone_id = ?", zId).Find(&postgresTrails).Error; err != nil {
		return nil, err
	}

	// Convert the postgresShelter records to domain.Shelter objects
	domainTrails := make([]*domain.Trail, len(postgresTrails))
	for i, ps := range postgresTrails {
		domainTrails[i] = toTrailAggregate(&ps)
	}
	return domainTrails, nil
}

// Shelters
func (repo *ShelterRepository) CreateShelter(id uuid.UUID, name string, tId uuid.UUID, availability bool, lat, long float64) (uuid.UUID, error) {
	shelter := postgresShelter{
		ShelterID:           id,
		ShelterName:         name,
		TrailID:             tId,
		ShelterAvailability: availability,
		Latitude:            lat,
		Longitude:           long,
	}
	if err := repo.db.Create(&shelter).Error; err != nil {
		return uuid.Nil, err
	}
	return shelter.ShelterID, nil
}

func (repo *ShelterRepository) UpdateShelterByID(id uuid.UUID, tId uuid.UUID, name string, availability bool, lat, long float64) error {
	return repo.db.Model(&postgresShelter{}).Where("shelter_id = ?", id).Updates(postgresShelter{
		ShelterName:         name,
		ShelterAvailability: availability,
		Latitude:            lat,
		Longitude:           long,
		TrailID:             tId,
	}).Error
}

func (repo *ShelterRepository) DeleteShelterByID(id uuid.UUID) error {
	return repo.db.Delete(&postgresShelter{}, "shelter_id = ?", id).Error
}

func (repo *ShelterRepository) GetShelterByID(id uuid.UUID) (*domain.Shelter, error) {
	var shelter postgresShelter
	if err := repo.db.Where("shelter_id = ?", id).First(&shelter).Error; err != nil {
		return nil, err // remove the domain.Shelter type
	}
	return toShelterAggregate(&shelter), nil
}

// GetAllShelters retrieves all shelter records from the database.
func (repo *ShelterRepository) List() ([]*domain.Shelter, error) {

	var postgresShelters []postgresShelter
	if err := repo.db.Find(&postgresShelters).Error; err != nil {
		return nil, err
	}

	// Convert the postgresShelter records to domain.Shelter objects
	domainShelters := make([]*domain.Shelter, len(postgresShelters))
	for i, ps := range postgresShelters {
		domainShelters[i] = toShelterAggregate(&ps)
	}

	return domainShelters, nil
}

func (repo *ShelterRepository) ListSheltersByTrailId(tId uuid.UUID) ([]*domain.Shelter, error) {
	var postgresShelters []postgresShelter
	if err := repo.db.Where("trail_id = ?", tId).Find(&postgresShelters).Error; err != nil {
		return nil, err
	}

	// Convert the postgresShelter records to domain.Shelter objects
	domainShelters := make([]*domain.Shelter, len(postgresShelters))
	for i, ps := range postgresShelters {
		domainShelters[i] = toShelterAggregate(&ps)

	}
	return domainShelters, nil
}

// functions for zone
func (repo *ZoneRepository) CreateZone(name string) (uuid.UUID, error) {
	zone := postgresZone{
		ZoneID:   uuid.New(),
		ZoneName: name,
	}
	if err := repo.db.Create(&zone).Error; err != nil {
		return uuid.Nil, err
	}
	return zone.ZoneID, nil
}

func (repo *ZoneRepository) UpdateZone(id uuid.UUID, name string) error {
	return repo.db.Model(&postgresZone{}).Where("zone_id = ?", id).Updates(postgresZone{
		ZoneID:   id,
		ZoneName: name,
	}).Error
}

func (repo *ZoneRepository) DeleteZone(id uuid.UUID) error {
	return repo.db.Delete(&postgresZone{}, "zone_id = ?", id).Error
}
