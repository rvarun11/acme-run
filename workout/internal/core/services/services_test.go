package services_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	amqpsecondaryadapter "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/amqp"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/clients"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"

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
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	trailID := uuid.New()
	HRMID := uuid.New()
	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, false)

	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil)
	peripheralClientMock.On("BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, true).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", workout.WorkoutID).Return(nil)

	// Test the Start function
	link, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)
	assert.Contains(t, link, "/workoutOptions?workoutID=")

	// Assert that the peripheral device was bound correctly
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, true)

	// Test the Stop function
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)

	if len(publisherMock.PublishedWorkouts) == 0 {
		t.Errorf("Expected PublishWorkoutStats to be called, but it wasn't")
	}

	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout)
	assert.True(t, stoppedWorkout.IsCompleted)
	assert.NotEmpty(t, stoppedWorkout.EndedAt)

	// Assert that the peripheral device was unbound correctly
	peripheralClientMock.AssertCalled(t, "UnbindPeripheralData", workout.WorkoutID)
}

/*
TestWorkoutService_StartWorkoutTwice:

	This test checks the behavior of the WorkoutService when trying to start a workout that is already active.
	It first starts a workout, then attempts to start the same workout again and expects an error.
	Finally, it stops the workout to check if the stop functionality works as expected.
*/
func TestWorkoutService_StartWorkoutTwice(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()
	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, false)

	// Mock expected calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	peripheralClientMock.On("BindPeripheralData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// Start the workout
	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Attempt to start the same workout again
	workout_two, _ := domain.NewWorkout(playerID, trailID, HRMID, false, false)
	_, secondStartErr := service.Start(&workout_two, HRMID, true)
	assert.Error(t, secondStartErr, "Expected an error when starting an already active workout")

	// Stop the workout
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}

/*
TestWorkoutService_HardcoreModeNoShelter:

	This test verifies that in hardcore mode, the shelter option (option 0) is not available.
	It starts a workout in hardcore mode, attempts to start the shelter option and expects an error,
	then stops the workout.
*/
func TestWorkoutService_HardcoreModeNoShelter(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()
	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, true) // Hardcore mode enabled

	// Mock expected calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	peripheralClientMock.On("BindPeripheralData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// Start the workout
	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Attempt to start shelter option in hardcore mode (expect error)
	shelterStartErr := service.StartWorkoutOption(workout.WorkoutID, 0) // Option 0 represents shelter
	assert.Error(t, shelterStartErr, "Expected an error when starting shelter option in hardcore mode")

	// Stop the workout
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}

/*
TestWorkoutService_WorkoutOptionsStartMultipleTimesStop:

	This test checks the behavior of the WorkoutService with regards to starting and stopping workout options.
	It involves starting a workout, then attempting various operations on workout options, including
	starting and stopping them under normal and erroneous conditions.
*/
func TestWorkoutService_WorkoutOptionsStartMultipleTimesStop(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()
	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, false)

	// Mock expected calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	peripheralClientMock.On("BindPeripheralData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// Start the workout
	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Attempt to stop a workout option when none have been started (expect error)
	stopOptionErr := service.StopWorkoutOption(workout.WorkoutID)
	assert.Error(t, stopOptionErr, "Expected an error when stopping a workout option that hasn't been started")

	// Start a workout option
	startOptionErr := service.StartWorkoutOption(workout.WorkoutID, 0)
	assert.NoError(t, startOptionErr)

	// Attempt to start the same workout option again (expect error)
	secondStartOptionErr := service.StartWorkoutOption(workout.WorkoutID, 1)
	assert.Error(t, secondStartOptionErr, "Expected an error when starting a workout option that's already active")

	// Stop the workout option
	stopOptionErr = service.StopWorkoutOption(workout.WorkoutID)
	assert.NoError(t, stopOptionErr)

	// Attempt to stop the workout option again (expect error)
	secondStopOptionErr := service.StopWorkoutOption(workout.WorkoutID)
	assert.Error(t, secondStopOptionErr, "Expected an error when stopping a workout option that's already stopped")

	// Stop the workout
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
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
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	trailID := uuid.New()
	HRMID := uuid.New()

	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, false)

	// Mock start workout with necessary steps
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil)

	// Mocked response for peripheral device client calls
	peripheralClientMock.On("BindPeripheralData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

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

		err := service.UpdateDistanceTravelled(workout.WorkoutID, lat, long, timeOfLocation)
		assert.NoError(t, err)

		expectedTotalDistance += 0.01 // Adding 10 meters for each update in km
		// there are 100 points
	}

	_, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)

	// Assert peripheral device unbind call
	peripheralClientMock.AssertCalled(t, "UnbindPeripheralData", workout.WorkoutID)

	// Get the total distance traveled
	actualTotalDistance, err := service.GetDistanceById(workout.WorkoutID)
	assert.NoError(t, err)

	// Assert the distance is as expected
	assert.InDelta(t, expectedTotalDistance*50000, actualTotalDistance, 0.01*50000, "The actual distance should be close to the expected distance")
}

/*
TestWorkoutProcess_Shelters:

	Basic Test to check the count of Shelters Taken
	The test asserts that the Shelters field in the workout data reflects the number of times
	the user has taken shelter during the workout session.
*/
func TestWorkoutProcess_Shelters(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()

	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, false)

	// Mocked responses for user service calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(false, nil) // Assuming hardcore mode affects shelter logic

	// Mocked response for peripheral device client calls
	peripheralClientMock.On("BindPeripheralData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// Start the workout using the service
	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Perform "Go to Shelter" action twice
	for i := 0; i < 2; i++ {
		err := service.StartWorkoutOption(workout.WorkoutID, services.ShelterBit)
		assert.NoError(t, err, "failed to start shelter option on iteration %d", i)

		err = service.StopWorkoutOption(workout.WorkoutID)
		assert.NoError(t, err, "failed to stop shelter option on iteration %d", i)
	}

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
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
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()

	// Create and start a workout instance
	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, true)

	// Mocked responses for user service calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	userClientMock.On("GetUserAge", playerID).Return(30, nil)              // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)
	randomHeartRate := uint8(rand.Intn(87) + 134)
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(randomHeartRate, nil)

	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false)

	// Get workout options and assert shelter is not an option
	links, err := service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	// In hardcore mode, shelter should not be present, verify it
	for _, link := range links {
		assert.NotContains(t, link.URL, "option=0", "Shelter option should not be present in hardcore mode")
	}

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
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
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()

	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, true)

	// Mocked responses for user service calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	userClientMock.On("GetUserAge", playerID).Return(30, nil)              // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	// First call, return a value less than 133
	firstHeartRate := uint8(rand.Intn(133)) // Random number between 0 and 132
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(firstHeartRate, nil).Once()

	// Second call, return a value greater than 133
	secondHeartRate := uint8(rand.Intn(87) + 134) // Random number between 134 and 255
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(secondHeartRate, nil).Once()

	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false)

	// Get workout options and assert shelter is not an option
	links, err := service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=2", "Escape must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=1", "Fight must go down")

	// Get workout options and assert shelter is not an option
	links, err = service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=1", "Fight must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=2", "Escape must go down")

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}

/*
TestWorkoutService_InitialWorkoutOptionsIfCardio:

	The cardio option which is Escape will have a higher ranking
	this will change when the player performs two escapes or more
*/
func TestWorkoutService_WorkoutOptionsIfCardio(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()

	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, true)

	// Mocked responses for user service calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("cardio", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	userClientMock.On("GetUserAge", playerID).Return(30, nil)              // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	randomHeartRate := uint8(rand.Intn(133))
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(randomHeartRate, nil)

	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false)

	// Get workout options and assert shelter is not an option
	links, err := service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=2", "Escape must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=1", "Fight must go down")

	err = service.StartWorkoutOption(workout.WorkoutID, 2)
	assert.NoError(t, err)
	err = service.StopWorkoutOption(workout.WorkoutID)
	assert.NoError(t, err)

	// Get workout options
	links, err = service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=2", "Escape must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=1", "Fight must go down")

	err = service.StartWorkoutOption(workout.WorkoutID, 2)
	assert.NoError(t, err)
	err = service.StopWorkoutOption(workout.WorkoutID)
	assert.NoError(t, err)

	// Get workout options
	links, err = service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	// Options should now be flipped
	assert.Contains(t, links[0].URL, "option=1", "Fight must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=2", "Escape must go down")

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
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
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()

	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, true)

	// Mocked responses for user service calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("strength", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	userClientMock.On("GetUserAge", playerID).Return(30, nil)              // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	randomHeartRate := uint8(rand.Intn(87) + 134)
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(randomHeartRate, nil)

	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false)

	// Get workout options
	links, err := service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	// Just check for fight followed by escape order
	assert.Contains(t, links[0].URL, "option=1", "Fight must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=2", "Escape must go down")

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}

/*
TestWorkoutService_InitialWorkoutOptionsIfStrength:

	The Strength option which is Fight will have a higher ranking
	this will change when the player performs two fights or more
*/
func TestWorkoutService_WorkoutOptionsIfStrength(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()

	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, true)

	// Mocked responses for user service calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("strength", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(true, nil) // Hardcore mode is on
	userClientMock.On("GetUserAge", playerID).Return(30, nil)              // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to false
	peripheralClientMock.On("BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	randomHeartRate := uint8(rand.Intn(87) + 134)
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(randomHeartRate, nil)

	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, false)

	// Get workout options
	links, err := service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=1", "Fight must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=2", "Escape must go down")

	err = service.StartWorkoutOption(workout.WorkoutID, 1)
	assert.NoError(t, err)
	err = service.StopWorkoutOption(workout.WorkoutID)
	assert.NoError(t, err)

	// Get workout options
	links, err = service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	assert.Contains(t, links[0].URL, "option=1", "Fight must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=2", "Escape must go down")

	err = service.StartWorkoutOption(workout.WorkoutID, 1)
	assert.NoError(t, err)
	err = service.StopWorkoutOption(workout.WorkoutID)
	assert.NoError(t, err)

	// Get workout options
	links, err = service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	// Options should now be flipped
	assert.Contains(t, links[0].URL, "option=2", "Escape must be at a higher rank")
	assert.Contains(t, links[1].URL, "option=1", "Fight must go down")

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}

/*
TestWorkoutService_DistanceToShelterUpdatesTest:

	Test to check that the incoming distance to shelter is updated in the workout options query
*/
func TestWorkoutService_DistanceToShelterUpdatesTest(t *testing.T) {
	// Initialize the mocks and the service
	userClientMock := clients.NewUserServiceClientMock()
	peripheralClientMock := clients.NewPeripheralDeviceClientMock()
	publisherMock := amqpsecondaryadapter.NewMockPublisher()
	store := postgres.NewRepository(cfg.Postgres)

	service := services.NewWorkoutService(store, peripheralClientMock, userClientMock, publisherMock)

	// Setup test data
	playerID := uuid.New()
	HRMID := uuid.New()
	trailID := uuid.New()

	workout, _ := domain.NewWorkout(playerID, trailID, HRMID, false, false)

	// Mocked responses for user service calls
	userClientMock.On("GetWorkoutPreferenceOfUser", playerID).Return("strength", nil)
	userClientMock.On("GetHardcoreModeOfUser", playerID).Return(false, nil) // Hardcore mode is off
	userClientMock.On("GetUserAge", playerID).Return(30, nil)               // Age or Player is 30

	// Mock the peripheral client to assert that the shelter request is set to true
	peripheralClientMock.On("BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, true).Return(nil)
	peripheralClientMock.On("UnbindPeripheralData", mock.Anything).Return(nil)

	randomHeartRate := uint8(rand.Intn(87) + 134)
	peripheralClientMock.On("GetAverageHeartRateOfUser", mock.Anything).Return(randomHeartRate, nil)

	_, startErr := service.Start(&workout, HRMID, true)
	assert.NoError(t, startErr)

	// Assert that the BindPeripheralData was called with shelterNeeded as false
	peripheralClientMock.AssertCalled(t, "BindPeripheralData", trailID, playerID, workout.WorkoutID, HRMID, true, true)

	// Mocking the Trail Manager
	distance := 10.0
	service.UpdateShelter(workout.WorkoutID, distance)

	// Get workout options and assert shelter is not an option
	links, err := service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	// In hardcore mode, shelter should not be present, verify it
	assert.Contains(t, links[0].Name, "Distance to Shelter = "+strconv.FormatFloat(distance, 'f', -1, 64)+" : ", "Distance Not Updated for Shelter")

	// Mocking the Trail Manager
	distance = 15.0
	service.UpdateShelter(workout.WorkoutID, distance)

	// Get workout options and assert shelter is not an option
	links, err = service.GetWorkoutOptions(workout.WorkoutID)
	assert.NoError(t, err)

	// In hardcore mode, shelter should not be present, verify it
	assert.Contains(t, links[0].Name, "Distance to Shelter = "+strconv.FormatFloat(distance, 'f', -1, 64)+" : ", "Distance Not Updated for Shelter")

	// Stop the workout using the service
	stoppedWorkout, stopErr := service.Stop(workout.WorkoutID)
	assert.NoError(t, stopErr)
	assert.NotNil(t, stoppedWorkout, "stopped workout should not be nil")
	assert.True(t, stoppedWorkout.IsCompleted, "stopped workout should be marked as completed")
}
