package aggregate_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
)

func TestWorkoutSession_NewWorkoutSession(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Empty workout session param validation",
			id:          uuid.Nil,
			expectedErr: aggregate.ErrInvalidWorkoutSession,
		}, {
			test:        "Valid workout session",
			id:          uuid.New(),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := aggregate.NewWorkoutSession(tc.id, false, true)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
