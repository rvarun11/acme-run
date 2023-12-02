package postgres

import (
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/domain"
	"github.com/google/uuid"
)

// Override the TableName method to specify the custom table name for the Zone model
func (postgresZone) TableName() string {
	return "zone"
}

// functions for zone
func (repo *DBRepository) CreateZone(name string) (uuid.UUID, error) {
	zone := postgresZone{
		ZoneID:   uuid.New(),
		ZoneName: name,
	}
	if err := repo.db.Create(&zone).Error; err != nil {
		return uuid.Nil, err
	}
	return zone.ZoneID, nil
}

func (repo *DBRepository) GetZoneByID(id uuid.UUID) (*domain.Zone, error) {
	var zone postgresZone
	if err := repo.db.Where("zone_id = ?", id).First(&zone).Error; err != nil {
		return nil, err // remove the domain.Shelter type
	}
	return zone.toAggregate(), nil
}

func (repo *DBRepository) GetZoneByName(name string) (*domain.Zone, error) {
	var zone postgresZone
	if err := repo.db.Where("zone_name = ?", name).First(&zone).Error; err != nil {
		return nil, err // remove the domain.Shelter type
	}
	return zone.toAggregate(), nil
}

func (repo *DBRepository) UpdateZone(id uuid.UUID, name string) error {
	return repo.db.Model(&postgresZone{}).Where("zone_id = ?", id).Updates(postgresZone{
		ZoneID:   id,
		ZoneName: name,
	}).Error
}

func (repo *DBRepository) DeleteZone(id uuid.UUID) error {
	return repo.db.Delete(&postgresZone{}, "zone_id = ?", id).Error
}

// List retrieves all zone records from the database.
func (repo *DBRepository) ListZones() ([]*domain.Zone, error) {
	var postgresZones []postgresZone
	if err := repo.db.Find(&postgresZones).Error; err != nil {
		return nil, err
	}

	// Convert the postgresZone records to domain.Zone objects
	domainZones := make([]*domain.Zone, len(postgresZones))
	for i, pzone := range postgresZones {
		domainZones[i] = pzone.toAggregate()
	}

	return domainZones, nil
}

func (repo *DBRepository) DeleteZoneByName(name string) error {
	return repo.db.Where("zone_name = ?", name).Delete(&postgresZone{}).Error
}
func toZonePostgres(zone *domain.Zone) *postgresZone {

	return &postgresZone{
		ZoneID:   zone.ZoneID,
		ZoneName: zone.ZoneName,
	}
}
