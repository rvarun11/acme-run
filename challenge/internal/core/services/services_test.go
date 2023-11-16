package services_test

import (
	"testing"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
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

var player1ID uuid.UUID = uuid.New()

type testStat struct {
	playerID        uuid.UUID
	distanceCovered float64
	fought          uint8
	escaped         uint8
	end             time.Time
}

func TestChallengeService_SubscribeToActiveChallenges(t *testing.T) {
	store := postgres.NewRepository(cfg.Postgres)
	service := services.NewChallengeService(store)
	// 1. 	Create Challenge
	chName := "Marathon Rush 2023" + uuid.NewString()
	ch, _ := domain.NewChallenge(chName, "", "", "DistanceCovered", 26.2, time.Now(), time.Now().Add(time.Second*30))
	ch, err := service.CreateChallenge(ch)
	if err != nil {
		logger.Error("something happened here", zap.Error(err))
	}

	// 2. Add Dummy DTOs
	one := testStat{
		playerID:        player1ID,
		distanceCovered: 10,
		fought:          2,
		escaped:         3,
		end:             time.Now().Add(time.Second),
	}
	two := testStat{
		playerID:        player1ID,
		distanceCovered: 20,
		fought:          2,
		escaped:         3,
		end:             time.Now().Add(time.Second * 2),
	}

	// 3. Subscribe to Challenge
	err = service.SubscribeToActiveChallenges(one.playerID, one.distanceCovered, one.fought, one.escaped, one.end)
	if err != nil {
		logger.Debug("unable to subscribe to challenge", zap.Error(err))
	}
	err = service.SubscribeToActiveChallenges(two.playerID, two.distanceCovered, two.fought, two.escaped, two.end)
	if err != nil {
		logger.Debug("unable to subscribe to challenge", zap.Error(err))
	}
	// Add delay here
	time.Sleep(15 * time.Second)

	// 3.  Get Badges
	badges, _ := service.ActeFinal(ch)
	logger.Debug("Badges List", zap.Any("badges", badges))
	// 4. Asert if the player received the badge or not
	t.Run("Badge Validation", func(t *testing.T) {
		if !badgeExists(badges, ch.ID, player1ID) {
			t.Errorf("badge not found")
		}
	})

}

func badgeExists(badges []*domain.Badge, challengeID, playerID uuid.UUID) bool {
	badgeExist := false
	for _, b := range badges {
		if b.PlayerID == playerID && b.Challenge.ID == challengeID {
			badgeExist = true
		}
	}
	return badgeExist
}
