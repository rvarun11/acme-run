package memory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
	workoutsession "github.com/rvarun11/macrun-teamvs/domain/workout_session"
)

type MemoryRepository struct {
	workoutSessions map[uuid.UUID]aggregate.WorkoutSession
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		workoutSessions: make(map[uuid.UUID]aggregate.WorkoutSession),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (aggregate.WorkoutSession, error) {
	if workoutSession, ok := mr.workoutSessions[id]; ok {
		return workoutSession, nil
	}
	return aggregate.WorkoutSession{}, workoutsession.ErrWorkoutSessionNotFound
}

func (mr *MemoryRepository) Add(ws aggregate.WorkoutSession) error {
	if mr.workoutSessions == nil {
		mr.Lock()
		mr.workoutSessions = make(map[uuid.UUID]aggregate.WorkoutSession)
		mr.Unlock()
	}
	if _, ok := mr.workoutSessions[ws.GetID()]; ok {
		return fmt.Errorf("workout session already exist: %w", workoutsession.ErrAddWorkoutSessionFailed)
	}
	mr.Lock()
	mr.workoutSessions[ws.GetID()] = ws
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(ws aggregate.WorkoutSession) error {
	if _, ok := mr.workoutSessions[ws.GetID()]; ok {
		return fmt.Errorf("workout session does not exist: %w", workoutsession.ErrorUpdateWorkoutSessionFailed)
	}
	mr.Lock()
	mr.workoutSessions[ws.GetID()] = ws
	mr.Unlock()
	return nil
}
