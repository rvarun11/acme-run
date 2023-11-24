package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
)

func TestChallenge_NewChallenge(t *testing.T) {
	type testCase struct {
		test        string
		expectedErr error
		name        string
		description string
		badgeURL    string
		criteria    string
		goal        float64
		start       time.Time
		end         time.Time
	}

	testCases := []testCase{
		{
			test:        "check invalid criteria",
			expectedErr: domain.ErrorInvalidCriteria,
			name:        "HalloweeK 2023",
			description: "Run 1000km in 1 week",
			badgeURL:    "https://www.something.com",
			criteria:    "Distance",
			goal:        1000.0,
			start:       time.Now(),
			end:         time.Now().Add(time.Hour),
		},
		{
			test:        "Invalid challenge duration",
			expectedErr: domain.ErrInvalidTime,
			name:        "HalloweeK 2023",
			description: "Run 1000km in 1 week",
			badgeURL:    "https://www.something.com",
			criteria:    "DistanceCovered",
			goal:        0.0,
			start:       time.Now().Add(time.Hour),
			end:         time.Now(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := domain.NewChallenge(tc.name, tc.description, tc.badgeURL, domain.Criteria(tc.criteria), tc.goal, tc.start, tc.end)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected err %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
