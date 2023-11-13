package tests

import (
	"math/rand"
	"testing"
	"time"

	logger "github.com/CAS735-F23/macrun-teamvsl/challenge_manager/log"
	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	"github.com/CAS735-F23/macrun-teamvsl/workout/tests/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var cfg *config.AppConfiguration = config.Config

/*
TestWorkoutService_StartAndStop:

	This test validates the "start" and "stop" workflow of a workout session within the WorkoutService.
	It checks whether a new workout can be initiated and properly terminated.
*/

func TestWorkoutService_StartAndStop(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := mocks.NewUserServiceClientMock()
	peripheralClientMock := mocks.NewPeripheralDeviceClientMock()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock)

	// Setup test data
	playerID := uuid.New()
	workoutID := uuid.New()
	HRMID := uuid.New()

	userClientMock.On("GetProfileOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil)
	peripheralClientMock.On("BindPeripheralData", playerID, workoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", workoutID).Return(nil)

	workout := &domain.Workout{
		WorkoutID: workoutID,
		PlayerID:  playerID,
		CreatedAt: time.Now(),
	}

	// Test the Start function
	link, startErr := service.Start(workout, HRMID, true)
	assert.NoError(t, startErr)
	assert.Contains(t, link, "/workoutOptions?workoutID=")

	// Assert that the peripheral device was bound correctly
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", playerID, workoutID, HRMID, true, false)

	// Test the Stop function
	stoppedWorkout, stopErr := service.Stop(workoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout)
	assert.True(t, stoppedWorkout.IsCompleted)
	assert.NotEmpty(t, stoppedWorkout.EndedAt)

	// Assert that the peripheral device was unbound correctly
	peripheralClientMock.AssertCalled(t, "UnbindPeripheralData", workoutID)
}

/*
TestWorkoutService_UpdateDistanceTravelled:

	This test case simulates the workout's distance tracking functionality by invoking
	UpdateDistanceTravelled multiple times with sequential latitude and longitude
	coordinates to simulate a user moving along a path. It ensures that distance
	calculations are accumulated correctly in the database. After simulating the workout path,
	it checks if the total distance reported by the GetDistanceById function matches the expected
	sum of distances from the updates.
*/
func TestWorkoutService_UpdateDistanceTravelled(t *testing.T) {
	// Mock setup
	userClientMock := mocks.NewUserServiceClientMock()
	peripheralClientMock := mocks.NewPeripheralDeviceClientMock()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock)

	// Setup test data
	playerID := uuid.New()
	workoutID := uuid.New()
	HRMID := uuid.New()

	// Mock start workout with necessary steps
	userClientMock.On("GetProfileOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil)
	peripheralClientMock.On("BindPeripheralData", playerID, workoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", workoutID).Return(nil)

	workout := domain.Workout{
		WorkoutID: workoutID,
		PlayerID:  playerID,
		CreatedAt: time.Now(),
	}

	// Assume the Start function initializes the workout correctly
	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Update distance traveled 100 times
	expectedTotalDistance := 0.0
	startLat, startLong := 40.730610, -73.935242
	endLat, endLong := 40.739604, -73.935242 // Approx 1000 mts north of the starting point

	for i := 0; i < 100; i++ {
		lat := startLat + float64(i)*(endLat-startLat)/100
		long := startLong + float64(i)*(endLong-startLong)/100
		timeOfLocation := time.Now().Add(time.Duration(rand.Intn(1000)) * time.Millisecond)

		err := service.UpdateDistanceTravelled(workoutID, lat, long, timeOfLocation)
		assert.NoError(t, err)

		expectedTotalDistance += 0.01 // Adding 10 meters for each update in km
		// there are 100 points
	}

	_, stopErr := service.Stop(workoutID)
	assert.NoError(t, stopErr)

	// Assert peripheral device unbind call
	peripheralClientMock.AssertCalled(t, "UnbindPeripheralData", workoutID)

	// Get the total distance traveled
	actualTotalDistance, err := service.GetDistanceById(workoutID)
	assert.NoError(t, err)

	// Assert the distance is as expected
	assert.InDelta(t, expectedTotalDistance, actualTotalDistance, 0.01, "The actual distance should be close to the expected distance")
}

/*
TestWorkoutProcess_Shelters:

	Basic Test to check the count of Shelters Taken
	The test asserts that the Shelters field in the workout data reflects the number of times
	the user has taken shelter during the workout session.
*/
func TestWorkoutProcess_Shelters(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := mocks.NewUserServiceClientMock()
	peripheralClientMock := mocks.NewPeripheralDeviceClientMock()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock)

	// Setup test data
	playerID := uuid.New()
	workoutID := uuid.New()
	HRMID := uuid.New()

	// Mocked responses for user service calls
	userClientMock.On("GetProfileOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(false, nil) // Assuming hardcore mode affects shelter logic

	// Mocked response for peripheral device client calls
	peripheralClientMock.On("BindPeripheralData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// Create a workout instance
	workout := &domain.Workout{
		WorkoutID: workoutID,
		PlayerID:  playerID,
		CreatedAt: time.Now(),
	}

	// Start the workout using the service
	_, startErr := service.Start(workout, HRMID, true)
	assert.NoError(t, startErr)

	// Perform "Go to Shelter" action twice
	for i := 0; i < 2; i++ {
		err := service.StartWorkoutOption(workoutID, services.ShelterBit)
		assert.NoError(t, err, "failed to start shelter option on iteration %d", i)

		err = service.StopWorkoutOption(workoutID)
		assert.NoError(t, err, "failed to stop shelter option on iteration %d", i)
	}

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")

	// Assert that the shelter count is as expected
	sheltersTaken := stoppedWorkout.Shelters
	assert.Equal(t, uint8(2), sheltersTaken, "shelters taken should be 2")
}

/*
TestWorkoutService_HardcoreMode:

	Test to check that Shelter is not given as an option for Hardcore mode
*/
func TestWorkoutService_HardcoreMode(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := mocks.NewUserServiceClientMock()
	peripheralClientMock := mocks.NewPeripheralDeviceClientMock()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock)

	// Setup test data
	playerID := uuid.New()
	workoutID := uuid.New()
	HRMID := uuid.New()

	// Mocked responses for user service calls
	userClientMock.On("GetProfileOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", playerID, workoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// Create and start a workout instance
	workout := &domain.Workout{
		WorkoutID:    workoutID,
		PlayerID:     playerID,
		CreatedAt:    time.Now(),
		HardcoreMode: true, // Set hardcore mode on the workout
	}

	_, startErr := service.Start(workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", playerID, workoutID, HRMID, true, false)

	// Get workout options and assert shelter is not an option
	links, err := service.GetWorkoutOptions(workoutID)
	assert.NoError(t, err)

	// In hardcore mode, shelter should not be present, verify it
	for _, link := range links {
		logger.Info(link.URL)
		assert.NotContains(t, link, "option=0", "Shelter option should not be present in hardcore mode")
	}

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}

/*
TestWorkoutService_InitialWorkoutOptionsIfCardio:

	Test to check that the options when the Profile is Cardio
	We flip the average heart rate to greater than 70% of the 220 - Age of User
	and then check the options, the options will also flip
	Meaning the Player should now prefer fighting than escaping
*/
func TestWorkoutService_InitialWorkoutOptionsIfCardio(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := mocks.NewUserServiceClientMock()
	peripheralClientMock := mocks.NewPeripheralDeviceClientMock()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock)

	// Setup test data
	playerID := uuid.New()
	workoutID := uuid.New()
	HRMID := uuid.New()

	// Mocked responses for user service calls
	userClientMock.On("GetProfileOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	userClientMock.On("GetUserAge", playerID).Return(30, nil)              // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", playerID, workoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// First call, return a value less than 133
	firstHeartRate := uint8(rand.Intn(133)) // Random number between 0 and 132
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(firstHeartRate, nil).Once()

	// Second call, return a value greater than 133
	secondHeartRate := uint8(rand.Intn(87) + 134) // Random number between 134 and 255
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(secondHeartRate, nil).Once()

	// Create and start a workout instance
	workout := &domain.Workout{
		WorkoutID:    workoutID,
		PlayerID:     playerID,
		CreatedAt:    time.Now(),
		HardcoreMode: true, // Set hardcore mode on the workout
	}

	_, startErr := service.Start(workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", playerID, workoutID, HRMID, true, false)

	service.ComputeWorkoutOptionsOrder(workoutID)

	// Get workout options and assert shelter is not an option
	links, err := service.GetWorkoutOptions(workoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=2", "Escape must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=1", "Fight must go down")

	service.ComputeWorkoutOptionsOrder(workoutID)

	// Get workout options and assert shelter is not an option
	links, err = service.GetWorkoutOptions(workoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=1", "Fight must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=2", "Escape must go down")

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}

/*
TestWorkoutService_InitialWorkoutOptionsIfStrength:

	Test to check that the options when the Profile is Strength
	It only checks the default options order which is fight followed by escape
*/
func TestWorkoutService_InitialWorkoutOptionsIfStrength(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := mocks.NewUserServiceClientMock()
	peripheralClientMock := mocks.NewPeripheralDeviceClientMock()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock)

	// Setup test data
	playerID := uuid.New()
	workoutID := uuid.New()
	HRMID := uuid.New()

	// Mocked responses for user service calls
	userClientMock.On("GetProfileOfUser", playerID).Return("strength", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	userClientMock.On("GetUserAge", playerID).Return(30, nil)              // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", playerID, workoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	randomHeartRate := uint8(rand.Intn(87) + 134)
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(randomHeartRate, nil)

	// Create and start a workout instance
	workout := &domain.Workout{
		WorkoutID:    workoutID,
		PlayerID:     playerID,
		CreatedAt:    time.Now(),
		HardcoreMode: true, // Set hardcore mode on the workout
	}

	_, startErr := service.Start(workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", playerID, workoutID, HRMID, true, false)

	// Get workout options and assert shelter is not an option
	links, err := service.GetWorkoutOptions(workoutID)
	assert.NoError(t, err)

	// In hardcore mode, shelter should not be present, verify it
	assert.Contains(t, links[0].URL, "option=1", "Fight must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=2", "Escape must go down")

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}
