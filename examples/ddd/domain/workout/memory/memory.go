package memory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
	workout "github.com/rvarun11/macrun-teamvs/domain/workout"
)

type MemoryRepository struct {
	Workouts map[uuid.UUID]aggregate.Workout
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		Workouts: make(map[uuid.UUID]aggregate.Workout),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (aggregate.Workout, error) {
	if Workout, ok := mr.Workouts[id]; ok {
		return Workout, nil
	}
	return aggregate.Workout{}, workout.ErrWorkoutNotFound
}

func (mr *MemoryRepository) Add(ws aggregate.Workout) error {
	if mr.Workouts == nil {
		mr.Lock()
		mr.Workouts = make(map[uuid.UUID]aggregate.Workout)
		mr.Unlock()
	}
	if _, ok := mr.Workouts[ws.GetID()]; ok {
		return fmt.Errorf("workout session already exist: %w", workout.ErrAddWorkoutFailed)
	}
	mr.Lock()
	mr.Workouts[ws.GetID()] = ws
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(ws aggregate.Workout) error {
	if _, ok := mr.Workouts[ws.GetID()]; ok {
		return fmt.Errorf("workout session does not exist: %w", workout.ErrorUpdateWorkoutFailed)
	}
	mr.Lock()
	mr.Workouts[ws.GetID()] = ws
	mr.Unlock()
	return nil
}
