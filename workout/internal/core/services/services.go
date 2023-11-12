package services

import (
	"fmt"
	"time"

	logger "github.com/CAS735-F23/macrun-teamvsl/challenge_manager/log"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/ports"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/umahmood/haversine"
)

type ActiveWorkoutsLastLocation struct {
	// Latitude of the Player
	Latitude float64 `json:"latitude"`
	// Longitude of the Player
	Longitude float64 `json:"longitude"`
	// Time of location
	TimeOfLocation time.Time `json:"time_of_location"`
}

type ActiveWorkoutsHeartRate struct {
	// HRM Connection
	HRMConnected bool
}

type WorkoutService struct {
	repo                       ports.WorkoutRepository
	peripheral                 ports.PeripheralDeviceClient
	user                       ports.UserServiceClient
	activeWorkoutsLastLocation map[uuid.UUID]ActiveWorkoutsLastLocation
	activeWorkoutsHeartRate    map[uuid.UUID]ActiveWorkoutsHeartRate
}

// Factory for creating a new WorkoutService
func NewWorkoutService(repo ports.WorkoutRepository, peripheral ports.PeripheralDeviceClient, user ports.UserServiceClient) *WorkoutService {
	return &WorkoutService{
		repo:       repo,
		peripheral: peripheral,
		user:       user,
	}
}

func (s *WorkoutService) List() ([]*domain.Workout, error) {
	return []*domain.Workout{}, nil
}

func (s *WorkoutService) GetWorkout(id uuid.UUID) (*domain.Workout, error) {
	return s.repo.GetWorkout(id)
}

func (s *WorkoutService) Start(workout *domain.Workout, HRMID uuid.UUID, HRMConnected bool) (string, error) {
	profile, err := s.user.GetProfileOfUser(workout.PlayerID)
	if err != nil {
		return "", err
	}

	hardcoreMode, err := s.user.GetHardcoreModeOfUser(workout.PlayerID)
	if err != nil {
		return "", err
	}
	workout.HardcoreMode = hardcoreMode
	workout.Profile = profile
	s.repo.UpdateWorkout(workout)

	s.peripheral.BindPeripheralData(workout.PlayerID, workout.WorkoutID, HRMID, HRMConnected)

	s.activeWorkoutsHeartRate[workout.WorkoutID] = ActiveWorkoutsHeartRate{
		HRMConnected: HRMConnected,
	}

	workoutOptions := &domain.WorkoutOptions{
		WorkoutID:             workout.WorkoutID,
		CurrentWorkoutOption:  7,
		FightsPushDown:        false,
		IsWorkoutOptionActive: false,
	}

	err = s.repo.Create(workout, workoutOptions)
	if err != nil {
		return "", err
	}
	logger.Info("Workout Created with ID : workoutID", zap.String("workoutID", workout.WorkoutID.String()))
	linkURL := fmt.Sprintf("/workoutOptions?workoutID=%s", workout.WorkoutID)
	return linkURL, err
}

func (s *WorkoutService) GetWorkoutOptions(workoutID uuid.UUID) (map[string]string, error) {
	pworkoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return nil, err
	}

	// Determine the order of workout options based on FightsPushDown
	optionsOrder := computeOptionsOrder(pworkoutOptions)

	// Generate HATEOAS links for StartWorkoutOption based on the computed order
	links := generateStartWorkoutOptionLinks(workoutID, optionsOrder)

	return links, nil
}

// Constants for bit positions
const (
	ShelterBit = 0
	FightBit   = 1
	EscapeBit  = 2
)

func getWorkoutType(bit int8) string {
	if bit == ShelterBit {
		return "Shelter"
	} else if bit == FightBit {
		return "Fight"
	}
	return "Escape"
}

// Compute the order of options based on FightsPushDown
func computeOptionsOrder(pworkoutOptions *domain.WorkoutOptions) []uint8 {
	// Initialize the default order
	order := []uint8{ShelterBit, FightBit, EscapeBit}

	// If FightsPushDown is set, change the order of Fight and Escape
	if pworkoutOptions.FightsPushDown {
		order[1], order[2] = order[2], order[1]
	}

	return order
}

// Generate HATEOAS links for StartWorkoutOption
func generateStartWorkoutOptionLinks(workoutID uuid.UUID, optionsOrder []uint8) map[string]string {
	links := make(map[string]string)

	for i, option := range optionsOrder {
		linkName := fmt.Sprintf("option%d", i+1)
		linkURL := fmt.Sprintf("/startWorkoutOption?workoutID=%s&option=%d", workoutID, option)
		links[linkName] = linkURL
	}

	return links
}

func (s *WorkoutService) UpdateDistanceTravelled(workoutID uuid.UUID, latitude float64, longitude float64, timeOfLocation time.Time) error {
	// Check if the workout ID exists in the location map
	lastLocation, locationExists := s.activeWorkoutsLastLocation[workoutID]

	if locationExists {
		// Calculate the distance between existing and new location
		distanceCovered := 0.0
		if lastLocation.Latitude != latitude || lastLocation.Longitude != longitude {
			// Calculate the distance covered using the Haversine formula
			// Create orb.Point for each coordinate
			point1 := haversine.Coord{Lat: lastLocation.Latitude, Lon: lastLocation.Longitude}
			point2 := haversine.Coord{Lat: latitude, Lon: longitude}

			// Calculate the distance using the Haversine formula
			_, distanceCovered = haversine.Distance(point1, point2)
		}

		// Update the workout distance if the distance covered is greater than 0
		if distanceCovered > 0 {
			// Get the workout from the repository
			workout, err := s.repo.GetWorkout(workoutID)
			if err != nil {
				return err // Propagate the error from the repository
			}

			// Update the workout distance
			workout.DistanceCovered += distanceCovered

			// Update the workout in the repository
			_, err = s.repo.UpdateWorkout(workout)
			if err != nil {
				return err // Propagate the error from the repository
			}
		}
	} else {
		// If the location doesn't exist, add it to the map
		workout, err := s.repo.GetWorkout(workoutID)
		if err != nil && !workout.IsCompleted {
			s.activeWorkoutsLastLocation[workoutID] = ActiveWorkoutsLastLocation{
				Latitude:       latitude,
				Longitude:      longitude,
				TimeOfLocation: timeOfLocation,
			}
		}
	}

	return nil // Return nil to indicate success
}

func (s *WorkoutService) UpdateShelter(workoutID uuid.UUID, DistanceToShelter float64) error {
	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	workoutOptions.DistanceToShelter = DistanceToShelter

	s.repo.UpdateWorkoutOptions(workoutOptions)

	if err != nil {
		return err // Propagate the error from the repository
	}

	return nil // Return nil to indicate success
}

func (s *WorkoutService) StartWorkoutOption(workoutID uuid.UUID, workoutType uint8) error {
	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Check if the workout option is already active
	if workoutOptions.IsWorkoutOptionActive {
		return domain.ErrWorkoutOptionAlreadyActive // Return an error with a custom message
	}

	// Update the workout option to make it active (you need to set appropriate fields)
	workoutOptions.IsWorkoutOptionActive = true
	workoutOptions.CurrentWorkoutOption = int8(workoutType)

	// Update the workout options in the repository
	_, err = s.repo.UpdateWorkoutOptions(workoutOptions)

	if err != nil {
		return err // Propagate the error from the repository
	}
	logger.Info("Workout Option OPT started for ID : workoutID", zap.String("workoutID", workoutOptions.WorkoutID.String()), zap.String("OPT", getWorkoutType(workoutOptions.CurrentWorkoutOption)))
	return nil // Return nil to indicate success
}

func (s *WorkoutService) StopWorkoutOption(workoutID uuid.UUID) error {
	// Get the workout options from the repository
	workout, err := s.repo.GetWorkout(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Check if the workout option is already inactive
	if !workoutOptions.IsWorkoutOptionActive {
		return domain.ErrWorkoutOptionAlreadyInActive
	}

	// Update the Shelters, Fights and Escapes
	if workoutOptions.CurrentWorkoutOption == ShelterBit {
		workout.Shelters++
	} else if workoutOptions.CurrentWorkoutOption == FightBit {
		workout.Fights++
	} else if workoutOptions.CurrentWorkoutOption == EscapeBit {
		workout.Escapes++
	}

	// Update the workout option to make it inactive (you need to set appropriate fields)
	workoutOptions.IsWorkoutOptionActive = false
	workoutOptions.CurrentWorkoutOption = -1

	// Update the workout options in the repository
	_, err = s.repo.UpdateWorkoutOptions(workoutOptions)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Update the workout in the repository
	_, err = s.repo.UpdateWorkout(workout)
	if err != nil {
		return err // Propagate the error from the repository
	}
	logger.Info("Workout Option OPT stopped for ID : workoutID", zap.String("workoutID", workoutOptions.WorkoutID.String()), zap.String("OPT", getWorkoutType(workoutOptions.CurrentWorkoutOption)))
	return nil // Return nil to indicate success
}

func (s *WorkoutService) Stop(id uuid.UUID) (*domain.Workout, error) {
	// Call Update() to update InProgress to False & EndedAt to time.Now()
	var tempWorkout *domain.Workout
	var err error
	tempWorkout, err = s.repo.GetWorkout(id)

	// TODO: Better error handling
	if err != nil {
		return nil, err
	}

	tempWorkout.EndedAt = time.Now()
	tempWorkout.IsCompleted = true

	s.repo.UpdateWorkout(tempWorkout)
	s.repo.DeleteWorkoutOptions(tempWorkout.WorkoutID)
	delete(s.activeWorkoutsLastLocation, tempWorkout.WorkoutID)
	delete(s.activeWorkoutsHeartRate, tempWorkout.WorkoutID)

	// Ask GeoHR to Stop
	s.peripheral.UnbindPeripheralData(tempWorkout.WorkoutID)

	// Notify ChallengeService
	return tempWorkout, err
}

func (s *WorkoutService) GetDistanceById(workoutID uuid.UUID) (float64, error) {
	return s.repo.GetDistanceByID(workoutID)
}

func (s *WorkoutService) GetDistanceCoveredBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (float64, error) {
	return s.repo.GetDistanceCoveredBetweenDates(playerID, startDate, endDate)
}

func (s *WorkoutService) GetEscapesMadeById(workoutID uuid.UUID) (uint16, error) {
	return s.repo.GetEscapesMadeByID(workoutID)
}

func (s *WorkoutService) GetEscapesMadeBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	return s.repo.GetEscapesMadeBetweenDates(playerID, startDate, endDate)
}

func (s *WorkoutService) GetFightsFoughtById(workoutID uuid.UUID) (uint16, error) {
	return s.repo.GetFightsFoughtByID(workoutID)
}

func (s *WorkoutService) GetFightsFoughtBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	return s.repo.GetFightsFoughtBetweenDates(playerID, startDate, endDate)
}

func (s *WorkoutService) GetSheltersTakenById(workoutID uuid.UUID) (uint16, error) {
	return s.repo.GetSheltersTakenByID(workoutID)
}

func (s *WorkoutService) GetSheltersTakenBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	return s.repo.GetSheltersTakenBetweenDates(playerID, startDate, endDate)
}

// Function to run periodically
func (s *WorkoutService) RunPeriodicTask() {
	for {
		// Iterate through active workouts in locationMap
		for workoutID := range s.activeWorkoutsLastLocation {
			_ = s.ComputeWorkoutOptionsOrder(workoutID)
			// TODO
			//if err != nil {
			// Handle the error (e.g., log it)
			//}
		}

		// Sleep for 4 seconds before the next iteration
		time.Sleep(4 * time.Second)
	}
}

// ComputeWorkoutOptionsOrder is modified to take profile directly
func (s *WorkoutService) ComputeWorkoutOptionsOrder(workoutID uuid.UUID) error {
	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Get the workout from the repository
	workout, err := s.repo.GetWorkout(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Calculate the weight based on cardio
	weight := calculateWeight(workout.Profile)

	// Get the number of fights and escapes
	fights, err := s.repo.GetFightsFoughtByID(workoutID)
	if err != nil {
		return err
	}
	escapes, err := s.repo.GetEscapesMadeByID(workoutID)
	if err != nil {
		return err
	}

	// Calculate the weight based on fights and escapes
	fightEscapeWeight := calculateFightEscapeWeight(fights, escapes)

	avgHeartRate, err := s.peripheral.GetAverageHeartRateOfUser(workout.WorkoutID)

	if err != nil {
		return err
	}

	// Calculate the weight based on average heart rate
	heartRateWeight := calculateHeartRateWeight(avgHeartRate)

	// Sum the weights
	totalWeight := weight + fightEscapeWeight + heartRateWeight

	// Update FightsPushDown based on the total weight
	if totalWeight >= 75 {
		workoutOptions.FightsPushDown = true
	} else {
		workoutOptions.FightsPushDown = false
	}

	// Update the workout options in the repository
	_, err = s.repo.UpdateWorkoutOptions(workoutOptions)
	if err != nil {
		return err // Propagate the error from the repository
	}

	return nil // Return nil to indicate success
}

// Helper function to calculate the weight based on cardio
func calculateWeight(profile string) int {
	if profile == "cardio" {
		return 50
	}
	return 0
}

// Helper function to calculate the weight based on fights and escapes
func calculateFightEscapeWeight(fights, escapes uint16) int {
	if fights-escapes >= 2 {
		return 25
	}
	return 0
}

// Helper function to calculate the weight based on average heart rate
func calculateHeartRateWeight(avgHeartRate uint8) int {
	// TODO: Fetch user age from profile or another service
	// age := profile.Age

	// Mock user age for testing
	age := 30

	maxHeartRate := 220 - age
	percentageMaxHeartRate := float64(avgHeartRate) / float64(maxHeartRate) * 100

	if percentageMaxHeartRate < 70 {
		return 25
	}

	return 0
}
