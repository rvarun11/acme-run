package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/ports"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/services"
	"github.com/google/uuid"
)

var cfg *config.AppConfiguration = config.Config

type testChallenge struct {
	// Challenge Details
	name     string
	desc     string
	badgeURL string
	criteria string
	goal     float64
}

type testStats struct {
	// Challenge Stats
	distanceCovered float64
	fought          uint8
	escaped         uint8
}

type testCase struct {
	test        string
	expectedErr error
	playerID    uuid.UUID
	challenge   *testChallenge
	stats       []*testStats
}

/*
This function checks subscribes the incoming challenges stats to active challenges.
Once the challenge ends, badge
Two types of challenges (DistancedCovered & EnemiesFoughtMoreThanEscape) are used for testing.
*/
func TestChallengeService_CreateOrUpdateChallengeStats(t *testing.T) {
	// 1. Test Setup
	store := postgres.NewRepository(cfg.Postgres)
	service := services.NewChallengeService(store)

	// 2. Test Scenario
	testCases := initTestCases()

	// 3. Run Tests
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// logger.Debug("Running new test case", zap.String("test", tc.test))
			// 2.1 Create a Challenge
			ch, err := domain.NewChallenge(tc.challenge.name+uuid.NewString(), tc.challenge.desc, tc.challenge.badgeURL, domain.Criteria(tc.challenge.criteria), tc.challenge.goal, time.Now(), time.Now().Add(time.Second*30))
			if err != nil {
				t.Errorf("unable to initialize challenge, got %v", err)
			}
			ch, err = service.CreateChallenge(ch)
			if err != nil {
				t.Errorf("unable to create the challenge, got %v", err)
			}

			// 2.2 Subscribe the Player to an active challenge with their stats
			for i, s := range tc.stats {
				err = service.CreateOrUpdateChallengeStats(tc.playerID, s.distanceCovered, s.fought, s.escaped, time.Now().Add(time.Second*time.Duration(i)))
				if err != nil {
					t.Errorf("unable to subscribe to active challenge, got %v", err)
				}
			}
			// added a delay for the challenge to end
			time.Sleep(40 * time.Second)

			// 2.3 Check for Badges
			badges, err := service.ListBadgesByPlayerID(tc.playerID)
			if err != nil {
				t.Errorf("unable to fetch badges, got %v", err)
			}
			// logger.Debug("Badges List", zap.Any("badges", badges))

			err = badgeExists(badges, ch.ID, tc.playerID)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected err %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func badgeExists(badges []*domain.Badge, challengeID, playerID uuid.UUID) error {
	for _, b := range badges {
		if b.PlayerID == playerID && b.Challenge.ID == challengeID {
			return nil
		}
	}
	return ports.ErrBadgeNotFound
}

func initTestCases() []*testCase {
	// Challenges
	marathonCh := &testChallenge{
		name:     "Marathon Rush",
		desc:     "",
		badgeURL: "",
		criteria: "DistanceCovered",
		goal:     26.2,
	}

	moreEnemiesFoughtCh := &testChallenge{
		name:     "HalloweeK",
		desc:     "",
		badgeURL: "",
		criteria: "FightMoreThanEscape",
		goal:     0.0,
	}

	testCases := []*testCase{
		{
			test:        "Completed distance covered challenge",
			expectedErr: nil,
			playerID:    uuid.New(),
			challenge:   marathonCh,
			stats: []*testStats{
				{
					distanceCovered: 10,
					fought:          0,
					escaped:         0,
				},
				{
					distanceCovered: 20,
					fought:          0,
					escaped:         0,
				},
			},
		},
		{
			test:        "Failed distance covered challenge",
			expectedErr: ports.ErrBadgeNotFound,
			playerID:    uuid.New(),
			challenge:   marathonCh,
			stats: []*testStats{
				{
					distanceCovered: 10,
					fought:          0,
					escaped:         0,
				},
				{
					distanceCovered: 5,
					fought:          0,
					escaped:         0,
				},
			},
		},
		{
			test:        "Completed more enemies fought than escaped challenge",
			expectedErr: nil,
			playerID:    uuid.New(),
			challenge:   moreEnemiesFoughtCh,
			stats: []*testStats{
				{
					distanceCovered: 5,
					fought:          1,
					escaped:         2,
				},
				{
					distanceCovered: 5,
					fought:          4,
					escaped:         1,
				},
			},
		},
		{
			test:        "Failed more enemies fought than escaped challenge",
			expectedErr: ports.ErrBadgeNotFound,
			playerID:    uuid.New(),
			challenge:   moreEnemiesFoughtCh,
			stats: []*testStats{
				{
					distanceCovered: 5,
					fought:          1,
					escaped:         2,
				},
				{
					distanceCovered: 5,
					fought:          1,
					escaped:         2,
				},
			},
		},
	}

	return testCases
}
