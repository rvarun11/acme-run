package services_test

import (
	"testing"

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
