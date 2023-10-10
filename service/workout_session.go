package service

import (
	"github.com/google/uuid"
	workoutsession "github.com/rvarun11/macrun-teamvs/domain/workout_session"
	"github.com/rvarun11/macrun-teamvs/domain/workout_session/memory"
)

type WorkoutConfiguration func(wss *WorkoutSessionService) error

type WorkoutSessionService struct {
	workoutSessions workoutsession.WorkoutSessionRepository
}

func NewWorkoutService(cfgs ...WorkoutConfiguration) (*WorkoutSessionService, error) {
	wss := &WorkoutSessionService{}

	for _, cfg := range cfgs {
		err := cfg(wss)

		if err != nil {
			return nil, err
		}
	}
	return wss, nil
}

// WithWorkoutSessionRepository applies a given WorkoutSession repository to the WorkoutService
func WithWorkoutSessionRepository(wsr workoutsession.WorkoutSessionRepository) WorkoutConfiguration {
	// return a function that matches the WorkoutConfiguration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(wss *WorkoutSessionService) error {
		wss.workoutSessions = wsr
		return nil
	}
}

// WithMemoryWorkoutRepository applies a memory WorkoutSession repository to the WorkoutService
func WithMemoryWorkoutSessionRepository() WorkoutConfiguration {
	// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
	wsr := memory.New()
	return WithWorkoutSessionRepository(wsr)
}

// CreateWorkoutSession will chaintogether all repositories to create a order for a customer
func (wss *WorkoutSessionService) CreateWorkoutSession(wsID uuid.UUID, productIDs []uuid.UUID) error {
	// Get the customer
	_, err := wss.workoutSessions.Get(wsID)
	if err != nil {
		return err
	}

	// Get each Product, Ouchie, We need a ProductRepository

	return nil
}
