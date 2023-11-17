package domain_test

import (
	"errors"
	"testing"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"
	"github.com/google/uuid"
)

func TestPlayer_NewPlayer(t *testing.T) {
	type testCase struct {
		test        string
		expectedErr error
		name        string
		email       string
		dob         string
		weight      float64
		height      float64
		pref        string
		zoneID      uuid.UUID
	}

	testCases := []testCase{
		{
			test:        "Empty weight validation",
			expectedErr: domain.ErrInvalidPlayerWeight,
			name:        "Percy Bolmer",
			email:       "percy@bolmer.com",
			dob:         "1998-19-08",
			weight:      0.0,
			height:      180.4,
			pref:        "cardio",
			zoneID:      uuid.New(),
		},
		{
			test:        "Empty height validation",
			expectedErr: domain.ErrInvalidPlayerHeight,
			name:        "Percy Bolmer",
			email:       "percy@bolmer.com",
			dob:         "1998-19-08",
			weight:      80.3,
			height:      0.0,
			pref:        "cardio",
			zoneID:      uuid.New(),
		},
		{
			test:        "Incorrect Preference validation",
			expectedErr: domain.ErrInvalidPlayerPreference,
			name:        "Percy Bolmer",
			email:       "percy@bolmer.com",
			dob:         "1998-19-08",
			weight:      80.3,
			height:      180.4,
			pref:        "jumping",
			zoneID:      uuid.New(),
		},
		{
			test:        "Empty zoneID validation",
			expectedErr: domain.ErrInvalidZoneID,
			name:        "Percy Bolmer",
			email:       "percy@bolmer.com",
			dob:         "1998-19-08",
			weight:      80.3,
			height:      180.4,
			pref:        "strength",
			zoneID:      uuid.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := domain.NewPlayer(tc.name, tc.email, tc.dob, tc.weight, tc.height, domain.Preference(tc.pref), tc.zoneID)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected err %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
