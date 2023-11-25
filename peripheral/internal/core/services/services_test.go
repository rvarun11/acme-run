package services_test

import (
	"fmt"
	"testing"
	"time"

	rabbitmqhandler "github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/handler/primary/rabbitmq"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/repository"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/secondary/clients"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

// var cfg *config.AppConfiguration = config.Config

/*
TestPeripheralService_Bind_Unbind:

	This test validates the "start" and "stop" workflow of a workout session within the WorkoutService.
	It checks whether a new workout can be initiated and properly terminated.
*/

func TestWorkoutService_StartAndStop(t *testing.T) {
	// Initialize the repository
	repo := repository.NewMemoryRepository()

	// Initialize the mocks
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()

	// Initialize the service with the repository and the mocks
	peripheralService := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	// Set up test data
	playerID := uuid.New()
	HRMID := uuid.New()
	workoutID := uuid.New()
	hrmConnect := true
	sendLiveLocation := true

	// Set expectations on the mocks
	// rabbitMQHandlerMock.On("SendLastLocation", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("time.Time")).Return(nil)
	// zoneClientMock.On("GetTrailLocation", mock.AnythingOfType("uuid.UUID")).Return(0.0, 0.0, 0.0, 0.0, nil)

	err := peripheralService.BindPeripheral(playerID, workoutID, HRMID, hrmConnect, sendLiveLocation)
	assert.NoError(t, err)

	// Optionally: You can check if the peripheral is now present in the memory repository
	pInstance, err := repo.GetByHRMId(HRMID)
	assert.NoError(t, err)
	assert.NotNil(t, pInstance)

	// Test unbinding a peripheral
	err = peripheralService.DisconnectPeripheral(workoutID)
	assert.NoError(t, err)

	// Optionally: You can check if the peripheral has been removed from the memory repository
	_, err = repo.GetByHRMId(HRMID)
	assert.Error(t, err)

}

func TestPeripheralService_UnbindWithoutBind(t *testing.T) {
	// Initialize the repository
	repo := repository.NewMemoryRepository()

	// Initialize the mocks without any expectations since they are not used here
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()

	// Initialize the service with the repository and the mocks
	peripheralService := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	// Set up test data
	workoutID := uuid.New()

	// Attempt to unbind a peripheral that was never bound
	err := peripheralService.DisconnectPeripheral(workoutID)

	// Assert that an error was returned
	assert.Error(t, err)

	// Optionally: Assert the type of error if your service returns different error types
	assert.Contains(t, err.Error(), ports.ErrorPeripheralNotFound.Error())
}

// TestDisconnectPeripheral_Exists checks unbinding an existing peripheral.
func TestDisconnectPeripheral_Exists(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	// First bind a peripheral
	pId := uuid.New()
	hId := uuid.New()
	wId := uuid.New()
	_ = service.BindPeripheral(pId, wId, hId, true, true)

	// Now unbind the peripheral
	err := service.DisconnectPeripheral(wId)
	assert.NoError(t, err)

	// Verify peripheral is no longer in the repository
	_, err = repo.GetByHRMId(hId)
	assert.Error(t, err)
}

// TestCreatePeripheral checks if a new peripheral is created correctly.
func TestCreatePeripheral(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	pId := uuid.New()
	hId := uuid.New()

	// Call the method under test
	err := service.CreatePeripheral(pId, hId)
	assert.NoError(t, err)

	// Further assertions to check if the peripheral is correctly added to the repository
	_, err = repo.GetByHRMId(hId)
	assert.NoError(t, err)
}

// TestCheckStatusByHRMId checks the status of a peripheral by HRM ID.
func TestCheckStatusByHRMId(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	// First create a peripheral
	pId := uuid.New()
	hId := uuid.New()
	_ = service.CreatePeripheral(pId, hId)

	// Now check the status
	status := service.CheckStatusByHRMId(hId)
	assert.False(t, status)

	// Check status for non-existing peripheral
	status = service.CheckStatusByHRMId(uuid.New())
	assert.False(t, status)
}

// TestBindPeripheral_NewPeripheral checks binding a new peripheral.
func TestBindPeripheral_NewPeripheral(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	pId := uuid.New()
	hId := uuid.New()
	wId := uuid.New()

	// Call the method under test
	err := service.BindPeripheral(pId, wId, hId, true, true)
	assert.NoError(t, err)

	// Verify peripheral is now in the repository and bound
	pInstance, err := repo.GetByHRMId(hId)
	assert.NoError(t, err)
	assert.Equal(t, pId, pInstance.PlayerId)
	assert.Equal(t, wId, pInstance.WorkoutId)
	assert.True(t, pInstance.HRMDev.HRMStatus)
	assert.True(t, pInstance.LiveStatus)
}

// TestDisconnectPeripheral_NotExists checks unbinding a peripheral that does not exist.
func TestDisconnectPeripheral_NotExists(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	wId := uuid.New()

	// Attempt to unbind a peripheral that was never bound
	err := service.DisconnectPeripheral(wId)
	// Verify that an error is returned since the peripheral does not exist
	assert.Error(t, err)
	// Optionally, verify the type of error if your service uses custom error types
	assert.Equal(t, ports.ErrorPeripheralNotFound, err)
}

// TestSetHeartRateReading tests updating the heart rate reading for a peripheral.
func TestSetHeartRateReading(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	// First, create and bind a peripheral
	pId := uuid.New()
	hId := uuid.New()
	wId := uuid.New()
	err := service.BindPeripheral(pId, wId, hId, true, true)
	assert.NoError(t, err)

	// Set a heart rate reading
	reading := 80 // Example heart rate reading
	err = service.SetHeartRateReading(hId, reading)
	assert.NoError(t, err)

	// Verify that the heart rate reading has been updated in the repository
	pInstance, err := repo.GetByHRMId(hId)
	assert.NoError(t, err)
	assert.Equal(t, reading, pInstance.HRMDev.HRate)
}

// TestGetHRMAvgReading checks retrieving the average heart rate reading.
func TestGetHRMAvgReading(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	// First, bind a peripheral with a known HRM ID and set an average heart rate
	pId := uuid.New()
	hId := uuid.New()
	wId := uuid.New()
	_ = service.BindPeripheral(pId, wId, hId, true, true)
	_ = service.SetHeartRateReading(hId, 80) // Example heart rate reading

	// Now get the average HRM reading
	hrmId, _, avgReading, err := service.GetHRMAvgReading(wId)
	fmt.Println(err)
	fmt.Println(hId, hrmId)
	assert.NoError(t, err)
	assert.Equal(t, hId, hrmId)
	assert.Equal(t, 80, avgReading) // Ensure the avgReading matches what was set
}

// TestGetHRMReading checks retrieving the current heart rate reading.
func TestGetHRMReading(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	pId := uuid.New()
	hId := uuid.New()
	wId := uuid.New()
	reading := 85
	_ = service.BindPeripheral(pId, wId, hId, true, true)
	_ = service.SetHeartRateReading(hId, reading)

	hrmId, timeRead, currentReading, err := service.GetHRMReading(hId)
	assert.NoError(t, err)
	assert.Equal(t, hId, hrmId)
	assert.Equal(t, reading, currentReading)
	assert.WithinDuration(t, time.Now(), timeRead, time.Second)
}

// TestGetHRMDevStatus checks retrieving the HRM device status for a peripheral.
func TestGetHRMDevStatus(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	hId := uuid.New()
	wId := uuid.New()
	_ = service.CreatePeripheral(wId, hId)
	status := true
	_ = service.SetHRMDevStatusByHRMId(hId, status)

	retrievedStatus, err := service.GetHRMDevStatus(wId)
	assert.NoError(t, err)
	assert.Equal(t, status, retrievedStatus)
}

// TestSetHRMDevStatusByHRMId tests setting the HRM device status by HRM ID.
func TestSetHRMDevStatusByHRMId(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	hId := uuid.New()
	pId := uuid.New()

	// Ensure that the peripheral is created successfully
	err := service.CreatePeripheral(pId, hId)
	assert.NoError(t, err)

	// Now try to update the status
	status := true
	err = service.SetHRMDevStatusByHRMId(hId, status)
	assert.NoError(t, err)

	// Retrieve the instance and check for a non-nil instance before asserting its status
	pInstance, err := repo.GetByHRMId(hId)
	assert.NoError(t, err)
	if pInstance != nil {
		assert.Equal(t, status, pInstance.HRMDev.HRMStatus)
	} else {
		t.Error("pInstance is nil")
	}
}

// TestSetGeoLocation tests setting the geolocation for a peripheral.
func TestSetGeoLocation(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	wId := uuid.New()
	hId := uuid.New()
	longitude := 40.712776
	latitude := -74.005974
	_ = service.CreatePeripheral(wId, hId)

	err := service.SetGeoLocation(wId, longitude, latitude)
	assert.NoError(t, err)

	pInstance, _ := repo.GetByWorkoutId(wId)
	assert.Equal(t, latitude, pInstance.GeoDev.Latitude)
	assert.Equal(t, longitude, pInstance.GeoDev.Longitude)
}

// TestGetGeoDevStatus checks retrieving the geolocation device status for a peripheral.
func TestGetGeoDevStatus(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	wId := uuid.New()
	hId := uuid.New()
	_ = service.CreatePeripheral(wId, hId)
	status := true
	_ = service.SetGeoDevStatus(wId, status)

	retrievedStatus, err := service.GetGeoDevStatus(wId)
	assert.NoError(t, err)
	assert.Equal(t, status, retrievedStatus)
}

// TestSetGeoDevStatus tests updating the geolocation device status for a peripheral.
func TestSetGeoDevStatus(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	wId := uuid.New()
	hId := uuid.New()
	_ = service.CreatePeripheral(wId, hId)
	status := true

	err := service.SetGeoDevStatus(wId, status)
	assert.NoError(t, err)

	pInstance, _ := repo.GetByWorkoutId(wId)
	assert.Equal(t, status, pInstance.GeoDev.GeoStatus)
}

// TestGetGeoLocation checks retrieving the current geolocation for a peripheral.
func TestGetGeoLocation(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	wId := uuid.New()
	hId := uuid.New()
	longitude := 40.712776
	latitude := -74.005974
	_ = service.CreatePeripheral(wId, hId)
	_ = service.SetGeoLocation(wId, longitude, latitude)

	_, retrievedLongitude, retrievedLatitude, _, err := service.GetGeoLocation(wId)
	assert.NoError(t, err)
	assert.Equal(t, latitude, retrievedLatitude)
	assert.Equal(t, longitude, retrievedLongitude)
}

// TestGetLiveStatus checks retrieving the live status for a peripheral.
func TestGetLiveStatus(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	wId := uuid.New()
	hId := uuid.New()
	_ = service.CreatePeripheral(wId, hId)
	liveStatus := true
	_ = service.SetLiveStatus(wId, liveStatus)

	retrievedStatus, err := service.GetLiveStatus(wId)
	assert.NoError(t, err)
	assert.Equal(t, liveStatus, retrievedStatus)
}

// TestSetLiveStatus tests setting the live status for a peripheral.
func TestSetLiveStatus(t *testing.T) {
	repo := repository.NewMemoryRepository()
	rabbitMQHandlerMock := rabbitmqhandler.NewRabbitMQHandlerMock()
	zoneClientMock := clients.NewZoneServiceClientMock()
	service := services.NewPeripheralService(repo, rabbitMQHandlerMock, zoneClientMock)

	wId := uuid.New()
	hId := uuid.New()
	_ = service.CreatePeripheral(wId, hId)
	liveStatus := true

	err := service.SetLiveStatus(wId, liveStatus)
	assert.NoError(t, err)

	pInstance, _ := repo.GetByWorkoutId(wId)
	assert.Equal(t, liveStatus, pInstance.LiveStatus)
}
