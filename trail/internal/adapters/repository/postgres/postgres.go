package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TrailRepository struct {
	db *gorm.DB
}

type ShelterRepo struct {
	db *gorm.DB
}

func NewTrailRepository(cfg *config.Postgres1) *Repository {

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
	db.AutoMigrate(&postgresWorkout{}, &postgresWorkoutOptions{})

	return &Repository{
		db: db,
	}
}

func NewShelterRepository(cfg *config.Postgres2) *Repository {

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
	db.AutoMigrate(&postgresWorkout{}, &postgresWorkoutOptions{})

	return &Repository{
		db: db,
	}
}

// Repository Types

type postgresTrail struct {
	TrailID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// name of the trail
	TrailName string `gorm:"type:string;unique not null"`
	// start longitude
	StartLongitude
	// start latitude
	StartLatitude
	// end longitude
	EndLongitude
	// end latitude
	EndLatitude
	// id of the cloest shelter
	CloestShelterId uuid.UUID
	// created time
	CreatedAt time.time
}

type postgresShelter struct {
	ShelterID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// name of the shelter
	ShelterName strig `gorm:"type:string;unique not null"`
	// availability of shelter
	ShelterAvailability
	// longitude of the shelter
	Longitude
	// latitude of the shelter
	Latitude
}

func toTrailAggregate(ptrail *postgresTrail) *domain.Trail {

	return &domain.Trail{
		TrailID:     ptrail.TrailID,
		TrailName:     ptrail.TrailName,
		StartLongitude:    ptrail.StartLongitude,
		StartLatitude: ptrail.StartLatitude,
		EndLongitude:   ptrail.EndLongitude,
		EndLatitude:     ptrail.EndLatitude,
		CloestShelterId: ptrail.CloestShelterId,
		CreatedAt:	ptrail.CreatedAt
	}
}

func toShelterAggregate(pshelter *postgresShelter) *domain.Shelter {

	return &domain.Shelter{
		ShelterID: pshelter.ShelterID
		ShelterName: pshelter.ShelterName
		ShelterAvailability: pshelter.ShelterAvailability
		Longitude: pshelter.Longitude
		Latitude: pshelter.Latitude
	}
}


func toTrailPostgres(trail *domain.Trail) *postgresTrail {

	return &postgresTrail{
		TrailID:     trail.TrailID,
		TrailName:     trail.TrailName,
		StartLongitude:    trail.StartLongitude,
		StartLatitude: trail.StartLatitude,
		EndLongitude:   trail.EndLongitude,
		EndLatitude:     trail.EndLatitude,
		CloestShelterId: trail.CloestShelterId,
		CreatedAt:	trail.CreatedAt
	}
}

func toShelterPostgres(shelter *domain.Shelter) *postgresShelter {

	return &postgresShelters{
		ShelterID: shelter.ShelterID
		ShelterName: shelter.ShelterName
		ShelterAvailability: shelter.ShelterAvailability
		Longitude: shelter.Longitude
		Latitude: shelter.Latitude
	}
}



// Repository Functions


func (repo *TrailRepository) CreateTrail(id uuid.UUID, name string, startLat, startLong, endLat, endLong float64, closestShelterID uuid.UUID) (uuid.UUID, error) {
	trail := postgresTrail{
		TrailID:          id,
		TrailName:        name,
		StartLatitude:    startLat,
		StartLongitude:   startLong,
		EndLatitude:      endLat,
		EndLongitude:     endLong,
		ClosestShelterId: closestShelterID,
		CreatedAt:        time.Now(),
	}
	if err := repo.DB.Create(&trail).Error; err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (repo *TrailRepository) UpdateTrailByID(id uuid.UUID, name string, startLat, startLong, endLat, endLong float64, closestShelterID uuid.UUID) error {
	return repo.DB.Model(&postgresTrail{}).Where("trail_id = ?", id).Updates(postgresTrail{
		TrailName:        name,
		StartLatitude:    startLat,
		StartLongitude:   startLong,
		EndLatitude:      endLat,
		EndLongitude:     endLong,
		ClosestShelterId: closestShelterID,
	}).Error
}

func (repo *TrailRepository) DeleteTrailByID(id uuid.UUID) error {
	return repo.DB.Delete(&postgresTrail{}, "trail_id = ?", id).Error
}


func (repo *TrailRepository) GetTrailByID(id uuid.UUID) (*domain.Trail, error) {
	var trail postgresTrail
	if err := repo.DB.Where("trail_id = ?", id).First(&trail).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return &domain.Trail, err
	}
	t := toTrailAggregate(trail)
	return t, nil
}

func (repo *ShelterRepository) CreateShelter(sid uuid.UUID, name string, availability bool, lat, long float64) (uuid.UUID, error) {
	shelter := postgresShelter{
		ShelterID:           sid,
		ShelterName:         name,
		ShelterAvailability: availability,
		Latitude:            lat,
		Longitude:           long,
	}
	if err := repo.DB.Create(&shelter).Error; err != nil {
		return uuid.Nil, err
	}
	return shelter.ShelterID, nil
}

func (repo *ShelterRepository) UpdateShelterByID(id uuid.UUID, name string, availability bool, lat, long float64) error {
	return repo.DB.Model(&postgresShelter{}).Where("shelter_id = ?", id).Updates(postgresShelter{
		ShelterName:         name,
		ShelterAvailability: availability,
		Latitude:            lat,
		Longitude:           long,
	}).Error
}

func (repo *ShelterRepository) DeleteShelterByID(id uuid.UUID) error {
	return repo.DB.Delete(&postgresShelter{}, "shelter_id = ?", id).Error
}

func (repo *ShelterRepository) GetShelterByID(id uuid.UUID) (*domain.Shelter, error) {
	var shelter postgresShelter
	if err := repo.DB.Where("shelter_id = ?", id).First(&shelter).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	s := toShelterAggregate(&shelter)
	return s, nil
}

// GetAllShelters retrieves all shelter records from the database.
func (repo *ShelterRepo) GetAllShelters() ([]*domain.Shelter, error) {
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



