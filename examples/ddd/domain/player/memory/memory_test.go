package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rvarun11/macrun-teamvs/aggregate"
	PlayerRepository "github.com/rvarun11/macrun-teamvs/domain/player"
)

func TestMemory_GetPlayer(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	// Create a fake customer to add to repository
	player, err := aggregate.NewPlayer("Samkith", "kishors@mcmaster.ca", "11/09/1997", 60.0, 5.6, "Mac")
	if err != nil {
		t.Fatal(err)
	}
	id := player.GetID()
	// Create the repo to use, and add some test Data to it for testing
	// Skip Factory for this
	repo := MemoryRepository{
		players: map[uuid.UUID]aggregate.Player{
			id: player,
		},
	}

	testCases := []testCase{
		{
			name:        "No Player By ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: PlayerRepository.ErrorPlayerNotFound,
		}, {
			name:        "Player By ID",
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := repo.Get(tc.id)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemory_AddPlayer(t *testing.T) {
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
			test:             "Add a Player",
			name:             "Samkith K Jain",
			email:            "kishors@mcmaster.ca",
			dob:              "11/09/1997",
			Weight:           60.0,
			Height:           5.6,
			GeographicalZone: "Mac",
			expectedErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := MemoryRepository{
				players: map[uuid.UUID]aggregate.Player{},
			}

			player, err := aggregate.NewPlayer(tc.name, tc.email, tc.dob, tc.Weight, tc.Height, tc.GeographicalZone)
			if err != nil {
				t.Fatal(err)
			}

			err = repo.Add(player)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

			found, err := repo.Get(player.GetID())
			if err != nil {
				t.Fatal(err)
			}
			if found.GetID() != player.GetID() {
				t.Errorf("Expected %v, got %v", player.GetID(), found.GetID())
			}
		})
	}
}
