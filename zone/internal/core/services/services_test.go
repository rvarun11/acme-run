package services_test

import (
	"testing"

	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	repository "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/memory"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/workouthandler"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var cfg *config.AppConfiguration = config.Config

func TestZoneService_CreateZoneManager(t *testing.T) {

	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()

	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	wId := uuid.New()
	zoneManagerID, err := service.CreateZoneManager(wId)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, zoneManagerID)

	// Additional logic to verify the ZoneManager in the database
	// ...
}

func TestZoneService_CreateTrail(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()

	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	trailName := "Test Trail"
	zoneID := uuid.New()
	startLatitude, startLongitude, endLatitude, endLongitude := 0.0, 0.0, 1.0, 1.0

	trailID, err := service.CreateTrail(trailName, zoneID, startLatitude, startLongitude, endLatitude, endLongitude)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, trailID)

}

func TestZoneService_CreateShelter(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()
	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	shelterName := "Test Shelter"
	trailID := uuid.New() // Assuming this trail already exists in your test setup
	availability := true
	lat, long := 40.7128, -74.0060

	shelterID, err := service.CreateShelter(shelterName, trailID, availability, lat, long)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, shelterID)

}

func TestZoneService_UpdateTrail(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()
	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	// Create a trail first
	trailName := "Original Trail Name"
	zoneID := uuid.New()
	trailID, _ := service.CreateTrail(trailName, zoneID, 0.0, 0.0, 1.0, 1.0)

	updatedName := "Updated Trail Name"
	err := service.UpdateTrail(trailID, updatedName, zoneID, 0.0, 0.0, 1.0, 1.0)

	assert.NoError(t, err)

	// Retrieve the updated trail and verify the changes
	updatedTrail, err := dbRepo.GetTrailByID(trailID)
	assert.NoError(t, err)
	assert.Equal(t, updatedName, updatedTrail.TrailName)
}

func TestZoneService_DeleteTrail(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()
	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	// Create a trail first
	trailName := "Test Trail"
	zoneID := uuid.New()
	trailID, _ := service.CreateTrail(trailName, zoneID, 0.0, 0.0, 1.0, 1.0)

	// Delete the trail
	err := service.DeleteTrail(trailID)
	assert.NoError(t, err)

	// Try to retrieve the deleted trail and expect an error
	_, err = dbRepo.GetTrailByID(trailID)
	assert.Error(t, err)
}

func TestZoneService_GetTrailByID(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()
	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	// Create a trail first
	trailName := "Test Trail"
	zoneID := uuid.New()
	trailID, _ := service.CreateTrail(trailName, zoneID, 0.0, 0.0, 1.0, 1.0)

	// Retrieve the trail
	retrievedTrail, err := service.GetTrailByID(trailID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTrail)
	assert.Equal(t, trailID, retrievedTrail.TrailID)
	assert.Equal(t, trailName, retrievedTrail.TrailName)
}
func TestZoneService_AddDuplicateTrail(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()
	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	trailName := "Unique Trail Name"
	zoneID := uuid.New()

	// Clean up before and after test
	defer dbRepo.DeleteTrailByName(trailName)
	dbRepo.DeleteTrailByName(trailName)

	// Add the trail for the first time
	_, err := service.CreateTrail(trailName, zoneID, 0.0, 0.0, 1.0, 1.0)
	assert.NoError(t, err)

	// Attempt to add the same trail again
	_, err = service.CreateTrail(trailName, zoneID, 0.0, 0.0, 1.0, 1.0)

	// Assert that an error is returned due to the duplicate name
	assert.Error(t, err, "Expected an error for duplicate trail creation, but got none")
}

func TestZoneService_AddDuplicateShelter(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()
	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	shelterName := "Unique Shelter Name"
	trailID := uuid.New() // Assuming this trail already exists

	// Clean up before and after test
	defer dbRepo.DeleteShelterByName(shelterName)
	dbRepo.DeleteShelterByName(shelterName)

	// Add the shelter for the first time
	_, err := service.CreateShelter(shelterName, trailID, true, 0.0, 0.0)
	assert.NoError(t, err)

	// Attempt to add the same shelter again
	_, err = service.CreateShelter(shelterName, trailID, true, 0.0, 0.0)

	// Assert that an error is returned due to the duplicate name
	assert.Error(t, err)
}

func TestZoneService_AddDuplicateZone(t *testing.T) {
	// Initialize repositories and service as above
	dbRepo := postgres.NewDBRepository(cfg.Postgres)
	zoneManagerRepo := repository.NewMemoryRepository()
	publisherMock := workouthandler.NewAMQPPublisherMock()
	service, _ := services.NewZoneService(zoneManagerRepo, dbRepo, publisherMock)

	zoneName := "Unique Zone Name"

	// Clean up before and after test
	defer dbRepo.DeleteZoneByName(zoneName)
	dbRepo.DeleteZoneByName(zoneName)

	// Add the zone for the first time
	_, err := service.CreateZone(zoneName)
	assert.NoError(t, err)

	// Attempt to add the same zone again
	_, err = service.CreateZone(zoneName)

	// Assert that an error is returned due to the duplicate name
	assert.Error(t, err)
}