package service

import (
	"github.com/google/uuid"
	workout "github.com/rvarun11/macrun-teamvs/domain/workout"
	"github.com/rvarun11/macrun-teamvs/domain/workout/memory"
)

type WorkoutConfiguration func(wss *WorkoutService) error

type WorkoutService struct {
	Workouts workout.WorkoutRepository
}

func NewWorkoutService(cfgs ...WorkoutConfiguration) (*WorkoutService, error) {
	wss := &WorkoutService{}

	for _, cfg := range cfgs {
		err := cfg(wss)

		if err != nil {
			return nil, err
		}
	}
	return wss, nil
}

// WithWorkoutRepository applies a given Workout repository to the WorkoutService
func WithWorkoutRepository(wsr workout.WorkoutRepository) WorkoutConfiguration {
	// return a function that matches the WorkoutConfiguration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(wss *WorkoutService) error {
		wss.Workouts = wsr
		return nil
	}
}

// WithMemoryWorkoutRepository applies a memory Workout repository to the WorkoutService
func WithMemoryWorkoutRepository() WorkoutConfiguration {
	// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
	wsr := memory.New()
	return WithWorkoutRepository(wsr)
}

// CreateWorkout will chaintogether all repositories to create a order for a customer
func (wss *WorkoutService) CreateWorkout(wsID uuid.UUID, productIDs []uuid.UUID) error {
	// Get the customer
	_, err := wss.Workouts.Get(wsID)
	if err != nil {
		return err
	}

	// Get each Product, Ouchie, We need a ProductRepository

	return nil
}
