package aggregate_test

import (
	"errors"
	"testing"

	"github.com/rvarun11/macrun-teamvs/aggregate"
)

func TestPlayer_NewPlayer(t *testing.T) {
	type testCase struct {
		test             string
		name             string
		email            string
		dob              string
		Weight           float32
		Height           float32
		GeographicalZone string
		expectedErr      error
	}

	testCases := []testCase{
		{
			test:             "Invalid Email Check",
			name:             "Samkith K Jain",
			email:            "kishors#mcmaster.ca",
			dob:              "11/09/1997",
			Weight:           60.0,
			Height:           60.0,
			GeographicalZone: "Mac",
			expectedErr:      aggregate.ErrInvalidEmail,
		},
		{
			test:             "Valid",
			name:             "Samkith K Jain",
			email:            "kishors@mcmaster.ca",
			dob:              "11/09/1997",
			Weight:           60.0,
			Height:           60.0,
			GeographicalZone: "Mac",
			expectedErr:      nil,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.test, func(t *testing.T) {
			_, err := aggregate.NewPlayer(tc.name, tc.email, tc.dob, tc.Weight, tc.Height, tc.GeographicalZone)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected Error %v, Got Error %v", tc.expectedErr, err)
			}
		})
	}
}
