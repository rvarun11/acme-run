package postgres

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/domain"
	"github.com/google/uuid"
)

// Override the TableName method to specify the custom table name for the Trail model
func (postgresTrail) TableName() string {
	return "trail"
}

// Repository Functions

func (repo *Repository) CreateTrail(name string, zId uuid.UUID, startLat, startLong, endLat, endLong float64) (uuid.UUID, error) {
	trail := postgresTrail{
		TrailID:        uuid.New(),
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
	return trail.TrailID, nil
}

func (repo *Repository) UpdateTrailByID(id uuid.UUID, name string, zId uuid.UUID, startLat, startLong, endLat, endLong float64) error {
	return repo.db.Model(&postgresTrail{}).Where("trail_id = ?", id).Updates(postgresTrail{
		TrailName:      name,
		ZoneID:         zId,
		StartLatitude:  startLat,
		StartLongitude: startLong,
		EndLatitude:    endLat,
		EndLongitude:   endLong,
	}).Error
}

func (repo *Repository) DeleteTrailByID(id uuid.UUID) error {
	return repo.db.Delete(&postgresTrail{}, "trail_id = ?", id).Error
}

func (repo *Repository) GetTrailByID(id uuid.UUID) (*domain.Trail, error) {
	var trail postgresTrail
	if err := repo.db.Where("trail_id = ?", id).First(&trail).Error; err != nil {
		return nil, err // remove the domain.Trail type
	}
	return trail.toAggregate(), nil
}

func (repo *Repository) ListTrails() ([]*domain.Trail, error) {
	var postgresTrails []postgresTrail
	if err := repo.db.Find(&postgresTrails).Error; err != nil {
		return nil, err
	}

	// Convert the postgresTrail records to domain.Trail objects
	domainTrails := make([]*domain.Trail, len(postgresTrails))
	for i, ptrail := range postgresTrails {
		domainTrails[i] = ptrail.toAggregate()
	}

	return domainTrails, nil
}
func (repo *Repository) ListTrailsByZoneId(zId uuid.UUID) ([]*domain.Trail, error) {
	var postgresTrails []postgresTrail
	if err := repo.db.Where("zone_id = ?", zId).Find(&postgresTrails).Error; err != nil {
		return nil, err
	}

	// Convert the postgresShelter records to domain.Shelter objects
	domainTrails := make([]*domain.Trail, len(postgresTrails))
	for i, ptrail := range postgresTrails {
		domainTrails[i] = ptrail.toAggregate()
	}
	return domainTrails, nil
}

func (repo *Repository) DeleteTrailByName(name string) error {
	return repo.db.Where("trail_name = ?", name).Delete(&postgresTrail{}).Error
}
