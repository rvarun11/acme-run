package domain_test

import (
	"errors"
	"testing"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"
)

func TestPlayer_NewUser(t *testing.T) {
	type testCase struct {
		test        string
		name        string
		email       string
		dob         string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Empty name validation",
			name:        "",
			expectedErr: domain.ErrInvalidUserName,
		},
		{
			test:        "Empty email validation",
			name:        "Percy Bolmer",
			email:       "",
			expectedErr: domain.ErrInvalidUserEmail,
		},
		{
			test:        "Empty dob validation",
			name:        "Percy Bolmer",
			email:       "percy@bolmer.com",
			dob:         "",
			expectedErr: domain.ErrInvalidUserDOB,
		},
		{
			test:        "Valid user",
			name:        "Percy Bolmer",
			email:       "percy@bolmer.com",
			dob:         "1998-19-08",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := domain.NewUser(tc.name, tc.email, tc.dob)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected err %v, got %v", tc.expectedErr, err)
			}
		})

	}
}
