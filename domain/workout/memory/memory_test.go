package memory

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
	workout "github.com/rvarun11/macrun-teamvs/domain/workout"
)

func TestMemory_GetWorkout(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	ws, err := aggregate.NewWorkout(uuid.New(), false, false)
	if err != nil {
		t.Fatal(err)
	}

	id := ws.GetID()

	repo := MemoryRepository{
		Workouts: map[uuid.UUID]aggregate.Workout{
			id: ws,
		},
	}

	testCases := []testCase{
		{
			name:        "no workout session by id",
			id:          uuid.MustParse("bd0776ac-581e-4a62-93d3-011ec4e072cd"),
			expectedErr: workout.ErrWorkoutNotFound,
		}, {
			name:        "workout session by id",
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.Get(tc.id)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}

}
