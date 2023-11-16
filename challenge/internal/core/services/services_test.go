package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/ports"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var cfg *config.AppConfiguration = config.Config

// Test Case 1: For Marathon Rush
// Test Setup to create challenge service
// Add Marathon Challenge
// Create Dummy CS DTOs - atleast 4 DTOs - two will get the badge, two won't
// Mock SubscribeTOACTIVeChs call with the above DTOs and give some delay (delay should be enough so that monitor challenge gets called)
// Get List of Badge and assert against the above players

// Test Case 2: Halloweek

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

func TestChallengeService_SubscribeToActiveChallenges(t *testing.T) {
	// 1. Test Setup
	store := postgres.NewRepository(cfg.Postgres)
	service := services.NewChallengeService(store)

	// 2. Test Scenario
	testCases := initTestCases()
	// 3. Run Tests
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			logger.Debug("Running new test case", zap.String("test", tc.test))
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
				err = service.SubscribeToActiveChallenges(tc.playerID, s.distanceCovered, s.fought, s.escaped, time.Now().Add(time.Second*time.Duration(i)))
				if err != nil {
					t.Errorf("unable to subscribe to active challenge, got %v", err)
				}
			}
			// added a delay for the challenge to end
			time.Sleep(15 * time.Second)

			// 2.3 Check for Badges
			badges, err := service.ActeFinal(ch)
			if err != nil {
				t.Errorf("expected err %v, got %v", tc.expectedErr, err)
			}
			logger.Debug("Badges List", zap.Any("badges", badges))

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
		criteria: "DistanceCovered",
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
					fought:          1,
					escaped:         2,
				},
			},
		},
	}

	return testCases
}

// // OLD
// func TestChallengeService_SubscribeToActiveChallenges_old(t *testing.T) {
// 	// 1. Test Setup
// 	store := postgres.NewRepository(cfg.Postgres)
// 	service := services.NewChallengeService(store)

// 	chName := "Marathon Rush 2023" + uuid.NewString()
// 	ch, _ := domain.NewChallenge(chName, "", "", "DistanceCovered", 26.2, time.Now(), time.Now().Add(time.Second*30))
// 	ch, err := service.CreateChallenge(ch)
// 	if err != nil {
// 		logger.Error("something happened here", zap.Error(err))
// 	}

// 	// 2. Add Dummy DTOs
// 	one := testStat{
// 		playerID:        player1ID,
// 		distanceCovered: 10,
// 		fought:          2,
// 		escaped:         3,
// 		end:             time.Now().Add(time.Second),
// 	}
// 	two := testStat{
// 		playerID:        player1ID,
// 		distanceCovered: 20,
// 		fought:          2,
// 		escaped:         3,
// 		end:             time.Now().Add(time.Second * 2),
// 	}

// 	// 3. Subscribe to Challenge
// 	err = service.SubscribeToActiveChallenges(one.playerID, one.distanceCovered, one.fought, one.escaped, one.end)
// 	if err != nil {
// 		logger.Debug("unable to subscribe to challenge", zap.Error(err))
// 	}
// 	err = service.SubscribeToActiveChallenges(two.playerID, two.distanceCovered, two.fought, two.escaped, two.end)
// 	if err != nil {
// 		logger.Debug("unable to subscribe to challenge", zap.Error(err))
// 	}
// 	// Add delay here
// 	time.Sleep(15 * time.Second)

// 	// 3.  Get Badges
// 	badges, _ := service.ActeFinal(ch)

// 	// 4. Asert if the player received the badge or not
// 	t.Run("Badge Validation", func(t *testing.T) {
// 		if !badgeExists(badges, ch.ID, player1ID) {
// 			t.Errorf("badge not found")
// 		}
// 	})

// }
