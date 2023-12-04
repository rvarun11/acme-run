package postgres

import (
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/domain"
	"github.com/google/uuid"
)

// Override the TableName method to specify the custom table name for the Shelter model
func (postgresShelter) TableName() string {
	return "shelter"
}

// Shelters
func (repo *Repository) CreateShelter(name string, tId uuid.UUID, availability bool, lat, long float64) (uuid.UUID, error) {
	shelter := postgresShelter{
		ShelterID:           uuid.New(),
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

func (repo *Repository) UpdateShelterByID(id uuid.UUID, tId uuid.UUID, name string, availability bool, lat, long float64) error {
	return repo.db.Model(&postgresShelter{}).Where("shelter_id = ?", id).Updates(postgresShelter{
		ShelterName:         name,
		ShelterAvailability: availability,
		Latitude:            lat,
		Longitude:           long,
		TrailID:             tId,
	}).Error
}

func (repo *Repository) DeleteShelterByID(id uuid.UUID) error {
	return repo.db.Delete(&postgresShelter{}, "shelter_id = ?", id).Error
}

func (repo *Repository) GetShelterByID(id uuid.UUID) (*domain.Shelter, error) {
	var shelter postgresShelter
	if err := repo.db.Where("shelter_id = ?", id).First(&shelter).Error; err != nil {
		return nil, err // remove the domain.Shelter type
	}
	return shelter.toAggregate(), nil
}

// GetAllShelters retrieves all shelter records from the database.
func (repo *Repository) ListShelters() ([]*domain.Shelter, error) {

	var postgresShelters []postgresShelter
	if err := repo.db.Find(&postgresShelters).Error; err != nil {
		return nil, err
	}

	// Convert the postgresShelter records to domain.Shelter objects
	domainShelters := make([]*domain.Shelter, len(postgresShelters))
	for i, shelter := range postgresShelters {
		domainShelters[i] = shelter.toAggregate()
	}

	return domainShelters, nil
}

func (repo *Repository) ListSheltersByTrailId(tId uuid.UUID) ([]*domain.Shelter, error) {
	var postgresShelters []postgresShelter
	if err := repo.db.Where("trail_id = ?", tId).Find(&postgresShelters).Error; err != nil {
		return nil, err
	}

	// Convert the postgresShelter records to domain.Shelter objects
	domainShelters := make([]*domain.Shelter, len(postgresShelters))
	for i, shelter := range postgresShelters {
		domainShelters[i] = shelter.toAggregate()

	}
	return domainShelters, nil
}

func (repo *Repository) DeleteShelterByName(name string) error {
	return repo.db.Where("shelter_name = ?", name).Delete(&postgresShelter{}).Error
}
